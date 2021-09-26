## An api to fetch covid cases from a state using goLang and echo framework.

**Endpoints**<br/>
**/**  -   GET request to check server is up and running and welcome message. <br/><br/>
**/fetchtodb** - GET request to fetch data from open covid api and persist it in mongodb database and display message after updating.<br/><br/>
**/state?latitude=value&longitude=value** -  GET request to fetch covid data for the state represented by given latitude and longitude and total cases in India.<br/>
                                                 Set latitude and longitude at the position of *value*.  <br/><br/>

This api is hosted at *https://nameless-reaches-06540.herokuapp.com/*.
