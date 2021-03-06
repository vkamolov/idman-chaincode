/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/

/*
Viktor: Identity Management Test
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
    "strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var cpPrefix      = "cp:"
var accountPrefix = "acct:"
var accountsKey   = "accounts"
/******* ID-Man *********************/
var personPrefix  = "pers:" 
var personKeysID  = "PersKeys"
var companyPrefix  = "comp:"
var companyKeysID  = "CompKeys"
/******* ID-Man *********************/

var recentLeapYear = 2016

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func generateCUSIPSuffix(issueDate string, days int) (string, error) {

	t, err := msToTime(issueDate)
	if err != nil {
		return "", err
	}

	maturityDate := t.AddDate(0, 0, days)
	month := int(maturityDate.Month())
	day := maturityDate.Day()

	suffix := seventhDigit[month] + eigthDigit[day]
	return suffix, nil

}

const (
	millisPerSecond     = int64(time.Second / time.Millisecond)
	nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt/millisPerSecond,
		(msInt%millisPerSecond)*nanosPerMillisecond), nil
}

/************* ID-Man **************************/
type UrlLink struct {
    Url         string   `json:"url"`
    UrlType     string   `json:"urlType"`
}

type Person struct {
	ID				string  `json:"id"`
	FirstName		string 	`json:"firstName"`
	LastName		string 	`json:"lastName"`
	Email			string 	`json:"email"`
	BirthDate		string 	`json:"birthDate"`
	Gender			string 	`json:"gender"`
	DrivingLicence	string 	`json:"drivingLicence"`
	TFN 			string 	`json:"tfn"`
	Address   		string  `json:"address"`
	City     		string  `json:"city"`
	Postcode 		string  `json:"postcode"`
	State    		string  `json:"state"`
	UrlLinks      []UrlLink `json:"urlLinks"`
	DataPhoto		string  `json:"dataPhoto"`
	Registrator    	string  `json:"registrator"`
	RegisterDate 	string  `json:"registerDate"`
}

type Company struct {
	ID				string  `json:"id"`
	Name			string 	`json:"name"`
	ACN 			string 	`json:"acn"`
	ABN 			string 	`json:"abn"`
	RegDate 		string 	`json:"regDate"`
	RegState		string 	`json:"regState"`
	Address   		string  `json:"address"`
	City     		string  `json:"city"`
	Postcode 		string  `json:"postcode"`
	State    		string  `json:"state"`
	UrlLinks      []UrlLink `json:"urlLinks"`
	Registrator    	string  `json:"registrator"`
	RegisterDate 	string  `json:"registerDate"`
}
/************* ID-Man **************************/


type Owner struct {
	Company string    `json:"company"`
	Quantity int      `json:"quantity"`
}

type CP struct {
	CUSIP     string  `json:"cusip"`
	Ticker    string  `json:"ticker"`
	Par       float64 `json:"par"`
	Qty       int     `json:"qty"`
	Discount  float64 `json:"discount"`
	Maturity  int     `json:"maturity"`
	Owners    []Owner `json:"owner"`
	Issuer    string  `json:"issuer"`
	IssueDate string  `json:"issueDate"`
}

type Account struct {
	ID          string  `json:"id"`
	Prefix      string  `json:"prefix"`
	CashBalance float64 `json:"cashBalance"`
	AssetsIds   []string `json:"assetIds"`
}

type Transaction struct {
	CUSIP       string   `json:"cusip"`
	FromCompany string   `json:"fromCompany"`
	ToCompany   string   `json:"toCompany"`
	Quantity    int      `json:"quantity"`
	Discount    float64  `json:"discount"`
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

/************* ID-Man **************************/
    // Initialize the collection of person keys
    fmt.Println("Initializing Person keys collection")
	
    var blank []string

    // Check if state already exists
    fmt.Println("Getting Person Keys")
    persKeysBytes, persErr := stub.GetState(personKeysID)
    if persKeysBytes == nil {
        fmt.Println("Cannot find PersKeys, will reinitialize everything")
        persBlankBytes, _ := json.Marshal(&blank)
        persErr := stub.PutState(personKeysID, persBlankBytes)
        if persErr != nil {
            fmt.Println("Failed to initialize person key collection")
        }
    } else if persErr != nil {
         fmt.Println("Failed to initialize person key collection")
    } else {
        fmt.Println("Found person keyBytes. Will not overwrite keys.")
    }

    // Initialize the collection of company keys
    fmt.Println("Initializing company keys collection")
	
    // Check if state already exists
    fmt.Println("Getting company Keys")
    compKeysBytes, compErr := stub.GetState(companyKeysID)
    if compKeysBytes == nil {
        fmt.Println("Cannot find company Keys, will reinitialize everything")
        compBlankBytes, _ := json.Marshal(&blank)
        compErr := stub.PutState(companyKeysID, compBlankBytes)
        if compErr != nil {
            fmt.Println("Failed to initialize company key collection")
        }
    } else if compErr != nil {
         fmt.Println("Failed to initialize company key collection")
    } else {
        fmt.Println("Found company keyBytes. Will not overwrite keys.")
    }
/************* ID-Man **************************/    
	
	fmt.Println("Initialization complete")
	return nil, nil
}

func (t *SimpleChaincode) createAccounts(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//  				0
	// "number of accounts to create"
	var err error
	numAccounts, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("error creating accounts with input")
		return nil, errors.New("createAccounts accepts a single integer argument")
	}
	//create a bunch of accounts
	var account Account
	counter := 1
	for counter <= numAccounts {
		var prefix string
		suffix := "000A"
		if counter < 10 {
			prefix = strconv.Itoa(counter) + "0" + suffix
		} else {
			prefix = strconv.Itoa(counter) + suffix
		}
		var assetIds []string
		account = Account{ID: "company" + strconv.Itoa(counter), Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
		accountBytes, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("error creating account" + account.ID)
			return nil, errors.New("Error creating account " + account.ID)
		}
		err = stub.PutState(accountPrefix+account.ID, accountBytes)
		counter++
		fmt.Println("created account" + accountPrefix + account.ID)
	}

	fmt.Println("Accounts created")
	return nil, nil

}

func (t *SimpleChaincode) createAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    // Obtain the username to associate with the account
    if len(args) != 1 {
        fmt.Println("Error obtaining username")
        return nil, errors.New("createAccount accepts a single username argument")
    }
    username := args[0]
    
    // Build an account object for the user
    var assetIds []string
    suffix := "000A"
    prefix := username + suffix
    var account = Account{ID: username, Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
    accountBytes, err := json.Marshal(&account)
    if err != nil {
        fmt.Println("error creating account" + account.ID)
        return nil, errors.New("Error creating account " + account.ID)
    }
    
    fmt.Println("Attempting to get state of any existing account for " + account.ID)
    existingBytes, err := stub.GetState(accountPrefix + account.ID)
	if err == nil {
        
        var company Account
        err = json.Unmarshal(existingBytes, &company)
        if err != nil {
            fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())
            
            if strings.Contains(err.Error(), "unexpected end") {
                fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
                err = stub.PutState(accountPrefix+account.ID, accountBytes)
                
                if err == nil {
                    fmt.Println("created account" + accountPrefix + account.ID)
                    return nil, nil
                } else {
                    fmt.Println("failed to create initialize account for " + account.ID)
                    return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
                }
            } else {
                return nil, errors.New("Error unmarshalling existing account " + account.ID)
            }
        } else {
            fmt.Println("Account already exists for " + account.ID + " " + company.ID)
		    return nil, errors.New("Can't reinitialize existing user " + account.ID)
        }
    } else {
        
        fmt.Println("No existing account found for " + account.ID + ", initializing account.")
        err = stub.PutState(accountPrefix+account.ID, accountBytes)
        
        if err == nil {
            fmt.Println("created account" + accountPrefix + account.ID)
            return nil, nil
        } else {
            fmt.Println("failed to create initialize account for " + account.ID)
            return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
        }
        
    }
}

/******* ID-Man *********************/
func (t *SimpleChaincode) registerPerson(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting person record")
	}

	var person Person
	var err error

	fmt.Println("Unmarshalling Person")
	err = json.Unmarshal([]byte(args[0]), &person)
	if err != nil {
		fmt.Println("error invalid person register")
		return nil, errors.New("Invalid Person register")
	}

	//generate the Person ID
	person.ID = strings.ToLower(person.FirstName) + strings.ToLower(person.LastName)
	person.ID = strings.Replace(person.ID, " ", "", -1) //remove all spaces
	//var stringHash := person.FirstName + person.LastName + person.BirthDate + person.Email + person.Gender
    //person.ID, err = genHash(stringHash)
    fmt.Println("Person ID is: ", person.ID)

    if person.ID == "" {
        fmt.Println("No Person ID, returning error")
        return nil, errors.New("Person ID cannot be blank")
    }
    fmt.Println("Person ID is: ", person.ID)
    fmt.Println("Person FirstName is: ", person.FirstName)
	fmt.Println("Person LastName is: ", person.LastName)
	fmt.Println("Person Email is: ", person.Email)
	fmt.Println("Person BirthDate is: ", person.BirthDate)
	fmt.Println("Person Gender is: ", person.Gender)
    fmt.Println("Person Address is: ", person.Address)
    fmt.Println("Person City is: ", person.City)
    fmt.Println("Person Postcode is: ", person.Postcode)
    fmt.Println("Person State is: ", person.State)
    fmt.Println("Registrator is: ", person.Registrator)
    fmt.Println("RegisterDate is: ", person.RegisterDate)

	fmt.Println("Marshalling Person bytes")
	fmt.Println("Getting State on Person " + person.ID)
	persRxBytes, err := stub.GetState(personPrefix+person.ID)

	if persRxBytes == nil {

		fmt.Println("ID does not exist, creating it")
		persBytes, err := json.Marshal(&person)
		if err != nil {
			fmt.Println("Error marshalling person")
			return nil, errors.New("Error registering person")
		}
		err = stub.PutState(personPrefix+person.ID, persBytes)
		if err != nil {
			fmt.Println("Error registering person")
			return nil, errors.New("Error registering person")
		}
		
		// Update the person keys by adding the new key
		fmt.Println("Getting Person Keys")
		keysBytes, err := stub.GetState(personKeysID)
		if err != nil {
			fmt.Println("Error retrieving person keys")
			return nil, errors.New("Error retrieving person keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling person keys ")
		}
		
		fmt.Println("Appending the new key to Person Keys")
		foundKey := false
		for _, key := range keys {
			if key == personPrefix+person.ID {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, personPrefix+person.ID)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on PersKeys")
			err = stub.PutState(personKeysID, keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}
		
		fmt.Println("Register person %+v\n", person)
		return nil, nil

	} else {

		
		fmt.Println("You can't create a person which already exists")
        return nil, errors.New("Can't a person which already exists")

        //Use for updating??
		fmt.Println("Person ID exists")
		
		var personRx Person
		fmt.Println("Unmarshalling Person " + person.ID)
		err = json.Unmarshal(persRxBytes, &personRx)
		if err != nil {
			fmt.Println("Error unmarshalling person " + person.ID)
			return nil, errors.New("Error unmarshalling person " + person.ID)
		}
		
		personRx.Address = person.Address
						
		persWriteBytes, err := json.Marshal(&personRx)
		if err != nil {
			fmt.Println("Error marshalling person")
			return nil, errors.New("Error registering a person")
		}
		err = stub.PutState(personPrefix+person.ID, persWriteBytes)
		if err != nil {
			fmt.Println("Error registering person")
			return nil, errors.New("Error registering person")
		}

		fmt.Println("Updated person %+v\n", personRx)
		return nil, nil
	}
}


func GetAllPersons(stub *shim.ChaincodeStub) ([]Person, error){
    
    var allPersons []Person
    
    // Get list of all the keys
    keysBytes, err := stub.GetState(personKeysID)
    if err != nil {
        fmt.Println("Error retrieving Person keys")
        return nil, errors.New("Error retrieving Person keys")
    }
    var keys []string
    err = json.Unmarshal(keysBytes, &keys)
    if err != nil {
        fmt.Println("Error unmarshalling Person keys")
        return nil, errors.New("Error unmarshalling Person keys")
    }

    // Get all the persons
    for _, value := range keys {
        persBytes, err := stub.GetState(value)
        
        var person Person
        err = json.Unmarshal(persBytes, &person)
        if err != nil {
            fmt.Println("Error retrieving person " + value)
            return nil, errors.New("Error retrieving person " + value)
        }
        
        fmt.Println("Appending Person" + value)
        allPersons = append(allPersons, person)
    }   
    
    return allPersons, nil
}

func GetPerson(personId string, stub *shim.ChaincodeStub) (Person, error){
    
    //
    persBytes, err := stub.GetState(personPrefix+personId)
    
    var person Person
    err = json.Unmarshal(persBytes, &person)
    if err != nil {
        fmt.Println("Error retrieving person " + personId)
        return person, errors.New("Error retrieving person " + personId)
    }
    
    return person, nil
}

func VerifyPerson(stub *shim.ChaincodeStub, sPerson string) (Person, error){

    var err error
    var person Person

    err = json.Unmarshal([]byte(sPerson), &person)

    if err != nil {
        return person, errors.New("Error unmarshalling verifying person")
    }

	//generate the person ID
	person.ID = strings.ToLower(person.FirstName) + strings.ToLower(person.LastName)
	person.ID = strings.Replace(person.ID, " ", "", -1) //remove all spaces
    fmt.Println("person ID is: ", person.ID)

    if person.ID == "" {
        fmt.Println("No person ID, returning error")
        return person, errors.New("person ID cannot be blank")
    }

    //Read existing person
    var personDB Person

	personDB, errDB := GetPerson(person.ID, stub)
	if errDB != nil {
		return person, errors.New("Person " + person.ID + " not found")
	}

	//Verifications (we don't check names. cause it's a part of the key)
	if 	(person.Email != personDB.Email) || (person.BirthDate != personDB.BirthDate) || (person.DrivingLicence != personDB.DrivingLicence) {

		return person, errors.New("Person verification failed")
	}

	return personDB, nil
}

func (t *SimpleChaincode) registerCompany(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting company record")
	}

	var company Company
	var err error

	fmt.Println("Unmarshalling company")
	err = json.Unmarshal([]byte(args[0]), &company)
	if err != nil {
		fmt.Println("error invalid company register")
		return nil, errors.New("Invalid company register")
	}

	//generate the company ID
	company.ID = strings.ToLower(company.Name)
	company.ID = strings.Replace(company.ID, " ", "", -1) //remove all spaces
	//var stringHash := person.FirstName + person.LastName + person.BirthDate + person.Email + person.Gender
    //person.ID, err = genHash(stringHash)
    fmt.Println("company ID is: ", company.ID)

    if company.ID == "" {
        fmt.Println("No company ID, returning error")
        return nil, errors.New("company ID cannot be blank")
    }
    fmt.Println("company ID is: ", company.ID)
    fmt.Println("company FirstName is: ", company.Name)
	fmt.Println("company ACN is: ", company.ACN)
	fmt.Println("company ABN is: ", company.ABN)
	fmt.Println("company RegDate is: ", company.RegDate)
	fmt.Println("company RegState is: ", company.RegState)
    fmt.Println("company Address is: ", company.Address)
    fmt.Println("company City is: ", company.City)
    fmt.Println("company Postcode is: ", company.Postcode)
    fmt.Println("company State is: ", company.State)
    fmt.Println("Registrator is: ", company.Registrator)
    fmt.Println("RegisterDate is: ", company.RegisterDate)

	fmt.Println("Marshalling company bytes")
	fmt.Println("Getting State on company " + company.ID)
	compRxBytes, err := stub.GetState(companyPrefix+company.ID)

	if compRxBytes == nil {

		fmt.Println("ID does not exist, creating it")
		compBytes, err := json.Marshal(&company)
		if err != nil {
			fmt.Println("Error marshalling company")
			return nil, errors.New("Error registering company")
		}
		err = stub.PutState(companyPrefix+company.ID, compBytes)
		if err != nil {
			fmt.Println("Error registering company")
			return nil, errors.New("Error registering company")
		}
		
		// Update the company keys by adding the new key
		fmt.Println("Getting company Keys")
		keysBytes, err := stub.GetState(companyKeysID)
		if err != nil {
			fmt.Println("Error retrieving company keys")
			return nil, errors.New("Error retrieving company keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling company keys ")
		}
		
		fmt.Println("Appending the new key to company Keys")
		foundKey := false
		for _, key := range keys {
			if key == companyPrefix+company.ID {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, companyPrefix+company.ID)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on company Keys")
			err = stub.PutState(companyKeysID, keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting company keys back")
				return nil, errors.New("Error writing the company keys back")
			}
		}
		
		fmt.Println("Register company %+v\n", company)
		return nil, nil

	} else {

		
		fmt.Println("You can't create a company which already exists")
        return nil, errors.New("Can't a company which already exists")

        //Use for updating??
		fmt.Println("company ID exists")
		
		var companyRx Company
		fmt.Println("Unmarshalling company " + company.ID)
		err = json.Unmarshal(compRxBytes, &companyRx)
		if err != nil {
			fmt.Println("Error unmarshalling company " + company.ID)
			return nil, errors.New("Error unmarshalling company " + company.ID)
		}
		
		companyRx.Address = company.Address
						
		compWriteBytes, err := json.Marshal(&companyRx)
		if err != nil {
			fmt.Println("Error marshalling company")
			return nil, errors.New("Error registering a company")
		}
		err = stub.PutState(companyPrefix+company.ID, compWriteBytes)
		if err != nil {
			fmt.Println("Error registering company")
			return nil, errors.New("Error registering company")
		}

		fmt.Println("Updated company %+v\n", companyRx)
		return nil, nil
	}
}


func GetAllCompanies(stub *shim.ChaincodeStub) ([]Company, error){
    
    var allCompanies []Company
    
    // Get list of all the keys
    keysBytes, err := stub.GetState(companyKeysID)
    if err != nil {
        fmt.Println("Error retrieving company keys")
        return nil, errors.New("Error retrieving company keys")
    }
    var keys []string
    err = json.Unmarshal(keysBytes, &keys)
    if err != nil {
        fmt.Println("Error unmarshalling company keys")
        return nil, errors.New("Error unmarshalling company keys")
    }

    // Get all the companies
    for _, value := range keys {
        compBytes, err := stub.GetState(value)
        
        var company Company
        err = json.Unmarshal(compBytes, &company)
        if err != nil {
            fmt.Println("Error retrieving company " + value)
            return nil, errors.New("Error retrieving company " + value)
        }
        
        fmt.Println("Appending company" + value)
        allCompanies = append(allCompanies, company)
    }   
    
    return allCompanies, nil
}

func GetCompany(companyId string, stub *shim.ChaincodeStub) (Company, error){
    
    //
    compBytes, err := stub.GetState(companyPrefix+companyId)
    
    var company Company
    err = json.Unmarshal(compBytes, &company)
    if err != nil {
        fmt.Println("Error retrieving company " + companyId)
        return company, errors.New("Error retrieving company " + companyId)
    }
    
    return company, nil
}

func VerifyCompany(stub *shim.ChaincodeStub, sCompany string) (Company, error){

    //sCompany = "{\"id\":\"test ltd\",\"name\":\"Test Ltd\"}"

    var err error
    var company Company

    err = json.Unmarshal([]byte(sCompany), &company)

    if err != nil {
        fmt.Println("Error retrieving company  + companyId")
        return company, errors.New("Error retrieving company  + companyId")
    }

	//generate the company ID
	company.ID = strings.ToLower(company.Name)
	company.ID = strings.Replace(company.ID, " ", "", -1) //remove all spaces
    fmt.Println("company ID is: ", company.ID)

    if company.ID == "" {
        fmt.Println("No company ID, returning error")
        return company, errors.New("company ID cannot be blank")
    }

    //Read existing company
    var companyDB Company

	companyDB, errDB := GetCompany(company.ID, stub)
	if errDB != nil {
		return company, errors.New("Company " + company.ID + " not found")
	}

	//Verifications (we don't check name. cause it's a part of the key)
	if 	(company.RegDate != companyDB.RegDate) || (company.RegState != companyDB.RegState) || (company.ACN != companyDB.ACN) || (company.ABN != companyDB.ABN) {

		return company, errors.New("Company verification failed")
	}

	return companyDB, nil
}

/******* ID-Man *********************/


func (t *SimpleChaincode) issueCommercialPaper(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	/*		0
		json
	  	{
			"ticker":  "string",
			"par": 0.00,
			"qty": 10,
			"discount": 7.5,
			"maturity": 30,
			"owners": [ // This one is not required
				{
					"company": "company1",
					"quantity": 5
				},
				{
					"company": "company3",
					"quantity": 3
				},
				{
					"company": "company4",
					"quantity": 2
				}
			],				
			"issuer":"company2",
			"issueDate":"1456161763790"  (current time in milliseconds as a string)

		}
	*/
	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
	}

	var cp CP
	var err error
	var account Account

	fmt.Println("Unmarshalling CP")
	err = json.Unmarshal([]byte(args[0]), &cp)
	if err != nil {
		fmt.Println("error invalid paper issue")
		return nil, errors.New("Invalid commercial paper issue")
	}

	//generate the CUSIP
	//get account prefix
	fmt.Println("Getting state of - " + accountPrefix + cp.Issuer)
	accountBytes, err := stub.GetState(accountPrefix + cp.Issuer)
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + cp.Issuer)
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		fmt.Println("Error Unmarshalling accountBytes")
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}
	
	account.AssetsIds = append(account.AssetsIds, cp.CUSIP)

	// Set the issuer to be the owner of all quantity
	var owner Owner
	owner.Company = cp.Issuer
	owner.Quantity = cp.Qty
	
	cp.Owners = append(cp.Owners, owner)

	suffix, err := generateCUSIPSuffix(cp.IssueDate, cp.Maturity)
	if err != nil {
		fmt.Println("Error generating cusip")
		return nil, errors.New("Error generating CUSIP")
	}

	fmt.Println("Marshalling CP bytes")
	cp.CUSIP = account.Prefix + suffix
	
	fmt.Println("Getting State on CP " + cp.CUSIP)
	cpRxBytes, err := stub.GetState(cpPrefix+cp.CUSIP)
	if cpRxBytes == nil {
		fmt.Println("CUSIP does not exist, creating it")
		cpBytes, err := json.Marshal(&cp)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(cpPrefix+cp.CUSIP, cpBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

		fmt.Println("Marshalling account bytes to write")
		accountBytesToWrite, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("Error marshalling account")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(accountPrefix + cp.Issuer, accountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on accountBytesToWrite")
			return nil, errors.New("Error issuing commercial paper")
		}
		
		
		// Update the paper keys by adding the new key
		fmt.Println("Getting Paper Keys")
		keysBytes, err := stub.GetState("PaperKeys")
		if err != nil {
			fmt.Println("Error retrieving paper keys")
			return nil, errors.New("Error retrieving paper keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling paper keys ")
		}
		
		fmt.Println("Appending the new key to Paper Keys")
		foundKey := false
		for _, key := range keys {
			if key == cpPrefix+cp.CUSIP {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, cpPrefix+cp.CUSIP)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on PaperKeys")
			err = stub.PutState("PaperKeys", keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}
		
		fmt.Println("Issue commercial paper %+v\n", cp)
		return nil, nil
	} else {
		fmt.Println("CUSIP exists")
		
		var cprx CP
		fmt.Println("Unmarshalling CP " + cp.CUSIP)
		err = json.Unmarshal(cpRxBytes, &cprx)
		if err != nil {
			fmt.Println("Error unmarshalling cp " + cp.CUSIP)
			return nil, errors.New("Error unmarshalling cp " + cp.CUSIP)
		}
		
		cprx.Qty = cprx.Qty + cp.Qty
		
		for key, val := range cprx.Owners {
			if val.Company == cp.Issuer {
				cprx.Owners[key].Quantity += cp.Qty
				break
			}
		}
				
		cpWriteBytes, err := json.Marshal(&cprx)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(cpPrefix+cp.CUSIP, cpWriteBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

		fmt.Println("Updated commercial paper %+v\n", cprx)
		return nil, nil
	}
}


func GetAllCPs(stub *shim.ChaincodeStub) ([]CP, error){
	
	var allCPs []CP
	
	// Get list of all the keys
	keysBytes, err := stub.GetState("PaperKeys")
	if err != nil {
		fmt.Println("Error retrieving paper keys")
		return nil, errors.New("Error retrieving paper keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling paper keys")
		return nil, errors.New("Error unmarshalling paper keys")
	}

	// Get all the cps
	for _, value := range keys {
		cpBytes, err := stub.GetState(value)
		
		var cp CP
		err = json.Unmarshal(cpBytes, &cp)
		if err != nil {
			fmt.Println("Error retrieving cp " + value)
			return nil, errors.New("Error retrieving cp " + value)
		}
		
		fmt.Println("Appending CP" + value)
		allCPs = append(allCPs, cp)
	}	
	
	return allCPs, nil
}

func GetCP(cpid string, stub *shim.ChaincodeStub) (CP, error){
	var cp CP

	cpBytes, err := stub.GetState(cpid)
	if err != nil {
		fmt.Println("Error retrieving cp " + cpid)
		return cp, errors.New("Error retrieving cp " + cpid)
	}
		
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + cpid)
		return cp, errors.New("Error unmarshalling cp " + cpid)
	}
		
	return cp, nil
}


/******* ID-Man cp-code *********************/
/*
func GetCompany(companyID string, stub *shim.ChaincodeStub) (Account, error){
	var company Account
	companyBytes, err := stub.GetState(accountPrefix+companyID)
	if err != nil {
		fmt.Println("Account not found " + companyID)
		return company, errors.New("Account not found " + companyID)
	}

	err = json.Unmarshal(companyBytes, &company)
	if err != nil {
		fmt.Println("Error unmarshalling account " + companyID + "\n err:" + err.Error())
		return company, errors.New("Error unmarshalling account " + companyID)
	}
	
	return company, nil
}
*/
/******* ID-Man cp-code *********************/

// Still working on this one
func (t *SimpleChaincode) transferPaper(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	/*		0
		json
	  	{
			  "CUSIP": "",
			  "fromCompany":"",
			  "toCompany":"",
			  "quantity": 1
		}
	*/
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting commercial paper record")
	}
	
	var tr Transaction

	fmt.Println("Unmarshalling Transaction")
	err := json.Unmarshal([]byte(args[0]), &tr)
	if err != nil {
		fmt.Println("Error Unmarshalling Transaction")
		return nil, errors.New("Invalid commercial paper issue")
	}

	fmt.Println("Getting State on CP " + tr.CUSIP)
	cpBytes, err := stub.GetState(cpPrefix+tr.CUSIP)
	if err != nil {
		fmt.Println("CUSIP not found")
		return nil, errors.New("CUSIP not found " + tr.CUSIP)
	}

	var cp CP
	fmt.Println("Unmarshalling CP " + tr.CUSIP)
	err = json.Unmarshal(cpBytes, &cp)
	if err != nil {
		fmt.Println("Error unmarshalling cp " + tr.CUSIP)
		return nil, errors.New("Error unmarshalling cp " + tr.CUSIP)
	}

	var fromCompany Account
	fmt.Println("Getting State on fromCompany " + tr.FromCompany)	
	fromCompanyBytes, err := stub.GetState(accountPrefix+tr.FromCompany)
	if err != nil {
		fmt.Println("Account not found " + tr.FromCompany)
		return nil, errors.New("Account not found " + tr.FromCompany)
	}

	fmt.Println("Unmarshalling FromCompany ")
	err = json.Unmarshal(fromCompanyBytes, &fromCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.FromCompany)
		return nil, errors.New("Error unmarshalling account " + tr.FromCompany)
	}

	var toCompany Account
	fmt.Println("Getting State on ToCompany " + tr.ToCompany)
	toCompanyBytes, err := stub.GetState(accountPrefix+tr.ToCompany)
	if err != nil {
		fmt.Println("Account not found " + tr.ToCompany)
		return nil, errors.New("Account not found " + tr.ToCompany)
	}

	fmt.Println("Unmarshalling tocompany")
	err = json.Unmarshal(toCompanyBytes, &toCompany)
	if err != nil {
		fmt.Println("Error unmarshalling account " + tr.ToCompany)
		return nil, errors.New("Error unmarshalling account " + tr.ToCompany)
	}

	// Check for all the possible errors
	ownerFound := false 
	quantity := 0
	for _, owner := range cp.Owners {
		if owner.Company == tr.FromCompany {
			ownerFound = true
			quantity = owner.Quantity
		}
	}
	
	// If fromCompany doesn't own this paper
	if ownerFound == false {
		fmt.Println("The company " + tr.FromCompany + "doesn't own any of this paper")
		return nil, errors.New("The company " + tr.FromCompany + "doesn't own any of this paper")	
	} else {
		fmt.Println("The FromCompany does own this paper")
	}
	
	// If fromCompany doesn't own enough quantity of this paper
	if quantity < tr.Quantity {
		fmt.Println("The company " + tr.FromCompany + "doesn't own enough of this paper")		
		return nil, errors.New("The company " + tr.FromCompany + "doesn't own enough of this paper")			
	} else {
		fmt.Println("The FromCompany owns enough of this paper")
	}
	
	amountToBeTransferred := float64(tr.Quantity) * cp.Par
	amountToBeTransferred -= (amountToBeTransferred) * (cp.Discount / 100.0) * (float64(cp.Maturity) / 360.0)
	
	// If toCompany doesn't have enough cash to buy the papers
	if toCompany.CashBalance < amountToBeTransferred {
		fmt.Println("The company " + tr.ToCompany + "doesn't have enough cash to purchase the papers")		
		return nil, errors.New("The company " + tr.ToCompany + "doesn't have enough cash to purchase the papers")	
	} else {
		fmt.Println("The ToCompany has enough money to be transferred for this paper")
	}
	
	toCompany.CashBalance -= amountToBeTransferred
	fromCompany.CashBalance += amountToBeTransferred

	toOwnerFound := false
	for key, owner := range cp.Owners {
		if owner.Company == tr.FromCompany {
			fmt.Println("Reducing Quantity from the FromCompany")
			cp.Owners[key].Quantity -= tr.Quantity
//			owner.Quantity -= tr.Quantity
		}
		if owner.Company == tr.ToCompany {
			fmt.Println("Increasing Quantity from the ToCompany")
			toOwnerFound = true
			cp.Owners[key].Quantity += tr.Quantity
//			owner.Quantity += tr.Quantity
		}
	}
	
	if toOwnerFound == false {
		var newOwner Owner
		fmt.Println("As ToOwner was not found, appending the owner to the CP")
		newOwner.Quantity = tr.Quantity
		newOwner.Company = tr.ToCompany
		cp.Owners = append(cp.Owners, newOwner)
	}
	
	fromCompany.AssetsIds = append(fromCompany.AssetsIds, tr.CUSIP)

	// Write everything back
	// To Company
	toCompanyBytesToWrite, err := json.Marshal(&toCompany)
	if err != nil {
		fmt.Println("Error marshalling the toCompany")
		return nil, errors.New("Error marshalling the toCompany")
	}
	fmt.Println("Put state on toCompany")
	err = stub.PutState(accountPrefix+tr.ToCompany, toCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the toCompany back")
		return nil, errors.New("Error writing the toCompany back")
	}
		
	// From company
	fromCompanyBytesToWrite, err := json.Marshal(&fromCompany)
	if err != nil {
		fmt.Println("Error marshalling the fromCompany")
		return nil, errors.New("Error marshalling the fromCompany")
	}
	fmt.Println("Put state on fromCompany")
	err = stub.PutState(accountPrefix+tr.FromCompany, fromCompanyBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the fromCompany back")
		return nil, errors.New("Error writing the fromCompany back")
	}
	
	// cp
	cpBytesToWrite, err := json.Marshal(&cp)
	if err != nil {
		fmt.Println("Error marshalling the cp")
		return nil, errors.New("Error marshalling the cp")
	}
	fmt.Println("Put state on CP")
	err = stub.PutState(cpPrefix+tr.CUSIP, cpBytesToWrite)
	if err != nil {
		fmt.Println("Error writing the cp back")
		return nil, errors.New("Error writing the cp back")
	}
	
	fmt.Println("Successfully completed Invoke")
	return nil, nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//need one arg
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ......")
	}

	if args[0] == "GetAllCPs" {
		fmt.Println("Getting all CPs")
		allCPs, err := GetAllCPs(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}	
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil		 
		}
	} else if args[0] == "GetCP" {
		fmt.Println("Getting particular cp")
		cp, err := GetCP(args[1], stub)
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			cpBytes, err1 := json.Marshal(&cp)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}	
			fmt.Println("All success, returning the cp")
			return cpBytes, nil		 
		}

/************* ID-Man **************************/	
	} else if args[0] == "GetAllPersons" {
		fmt.Println("Getting all Persons")
		allPersons, err := GetAllPersons(stub)
		if err != nil {
			fmt.Println("Error from GetAllPersons")
			return nil, err
		} else {
			allPersonsBytes, err1 := json.Marshal(&allPersons)
			if err1 != nil {
				fmt.Println("Error marshalling allPersons")
				return nil, err1
			}	
			fmt.Println("All success, returning allPersons")
			return allPersonsBytes, nil		 
		}

	} else if args[0] == "GetPerson" {
		fmt.Println("Getting particular person")
		person, err := GetPerson(args[1], stub)
		if err != nil {
			fmt.Println("Error Getting particular person")
			return nil, err
		} else {
			personBytes, err1 := json.Marshal(&person)
			if err1 != nil {
				fmt.Println("Error marshalling the person")
				return nil, err1
			}	
			fmt.Println("All success, returning the person")
			return personBytes, nil		 
		}

	} else if args[0] == "GetAllCompanies" {
		fmt.Println("Getting all Companies")
		allCompanies, err := GetAllCompanies(stub)
		if err != nil {
			fmt.Println("Error from GetAllCompanies")
			return nil, err
		} else {
			allCompaniesBytes, err1 := json.Marshal(&allCompanies)
			if err1 != nil {
				fmt.Println("Error marshalling allCompanies")
				return nil, err1
			}	
			fmt.Println("All success, returning allCompanies")
			return allCompaniesBytes, nil		 
		}		

	} else if args[0] == "GetCompany" {
		fmt.Println("Getting the company")
		company, err := GetCompany(args[1], stub)
		if err != nil {
			fmt.Println("Error from getCompany")
			return nil, err
		} else {
			companyBytes, err1 := json.Marshal(&company)
			if err1 != nil {
				fmt.Println("Error marshalling the company")
				return nil, err1
			}	
			fmt.Println("All success, returning the company")
			return companyBytes, nil		 
		}

	} else if args[0] == "VerifyCompany" {
		fmt.Println("Verifying the company")
		company, err := VerifyCompany(stub, args[1])
		if err != nil {
			fmt.Println("Error from VerifyCompany")
			return nil, err
		} else {
			companyBytes, err1 := json.Marshal(&company)
			if err1 != nil {
				fmt.Println("Error marshalling the company")
				return nil, err1
			}	
			fmt.Println("All success, returning the company")
			return companyBytes, nil		 
		}

	} else if args[0] == "VerifyPerson" {
		fmt.Println("Verifying the person")
		person, err := VerifyPerson(stub, args[1])
		if err != nil {
			fmt.Println("Error from VerifyPerson")
			return nil, err
		} else {
			personBytes, err1 := json.Marshal(&person)
			if err1 != nil {
				fmt.Println("Error marshalling the person")
				return nil, err1
			}	
			fmt.Println("All success, returning the person")
			return personBytes, nil		 
		}	

/************* ID-Man **************************/

	} else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0])

		if err != nil {
			fmt.Println("Some error happenend")
			return nil, errors.New("Some Error happened")
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil		
	} 
}

func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	
	if function == "issueCommercialPaper" {
		fmt.Println("Firing issueCommercialPaper")
		//Create an asset with some value
		return t.issueCommercialPaper(stub, args)

/************* ID-Man **************************/		
	} else if function == "registerPerson" {
        //Create a Person
        return t.registerPerson(stub, args)	

	} else if function == "registerCompany" {
        //Create a Company
        return t.registerCompany(stub, args)
/************* ID-Man **************************/      
	} else if function == "transferPaper" {
		fmt.Println("Firing cretransferPaperateAccounts")
		return t.transferPaper(stub, args)	
	} else if function == "createAccounts" {
		fmt.Println("Firing createAccounts")
		return t.createAccounts(stub, args)
	} else if function == "createAccount" {
        fmt.Println("Firing createAccount")
        return t.createAccount(stub, args)
    } else if function == "init" {
        fmt.Println("Firing init")
        return t.Init(stub, "init", args)
    }

	return nil, errors.New("Received unknown function invocation")
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}

//lookup tables for last two digits of CUSIP
var seventhDigit = map[int]string{
	1:  "A",
	2:  "B",
	3:  "C",
	4:  "D",
	5:  "E",
	6:  "F",
	7:  "G",
	8:  "H",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "P",
	15: "Q",
	16: "R",
	17: "S",
	18: "T",
	19: "U",
	20: "V",
	21: "W",
	22: "X",
	23: "Y",
	24: "Z",
}

var eigthDigit = map[int]string{
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "A",
	11: "B",
	12: "C",
	13: "D",
	14: "E",
	15: "F",
	16: "G",
	17: "H",
	18: "J",
	19: "K",
	20: "L",
	21: "M",
	22: "N",
	23: "P",
	24: "Q",
	25: "R",
	26: "S",
	27: "T",
	28: "U",
	29: "V",
	30: "W",
	31: "X",
}
