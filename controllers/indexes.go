package controllers

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"loomio/config"
)

func AddSecondaryIndex(dbSession config.Database, userid string, email string) error {

	/*
	   TO create a secondary index the best apprach would be a B-Tree, but in
	   here A new static key will be stored in LevelDB with value as a map[string]string
	   in the format
	       "email" -> userid


	*/
	secondaryIndex := map[string]string{}
	data, err := dbSession.DBSession.Get([]byte("secondaryIndex"), nil)
	if err != nil {
		fmt.Println("Empty secondary indexes", err)
	}

	reader := bytes.NewReader(data)
	// make a decoder
	decoder := gob.NewDecoder(reader)
	// decode it int

	decoder.Decode(&secondaryIndex)
	log.Println("This is hte secondasyindex data", secondaryIndex)
	secondaryIndex[email] = userid

	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)

	err = enc.Encode(secondaryIndex)
	if err != nil {
		log.Println("Error in encoding gob")
		return err
	}

	err = dbSession.DBSession.Put([]byte("secondaryIndex"), network.Bytes(), nil)
	//dberr := userCollection.Insert(user)
	if err != nil {
		log.Println("Error in Updatinf secondary indexes", err)
		return err
	}

	return nil
}

func RetreiveSecondaryIndex(dbSession config.Database, email string) bool {
	secondaryIndex := map[string]string{}
	data, _ := dbSession.DBSession.Get([]byte("secondaryIndex"), nil)
	reader := bytes.NewReader(data)
	// make a decoder
	decoder := gob.NewDecoder(reader)
	// decode it int

	decoder.Decode(&secondaryIndex)

	userid, ok := secondaryIndex[email]
	log.Println("This is the email to be found in secondary index", userid)

	return ok
}
