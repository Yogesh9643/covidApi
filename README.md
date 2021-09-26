## An api to fetch covid cases from a state using goLang and echo framework.

**Endpoints**
**/**  -   GET request to check data server is up and running and welcome message.
**/fetchtodb** - GET request to fetch data from open covid api and persist it in mongodb database.
**/state?longitude=<value>&latitude=<value>** -  GET request to fetch covid data for the state represented by given longitude and latitude and total cases in india.
                                                 Set longitude and latitude at the position of <value>

This api is hosted at *https://nameless-reaches-06540.herokuapp.com/*
