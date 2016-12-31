
boltapi wraps <a href="https://github.com/boltdb/bolt">boltdb</a> with a REST API interface. 
It allows the user to easily perform CRUD operations for key/value pairs as well as create/delete entire databases.


*********************************************************************************************
PUT OPERATIONS:

Databases & buckets are created (if it does not already exist) when a key/value pair is written using PUT.

/dbs/{db}/bucket/{bucketName}/keys/{keyName}
Creates Database and/or bucket if they do not already exist. Inserts key/value pair using BODY as value.

    Ex:
            PUT 	>> 	http://host:8080/dbs/somenewdatabase.db/buckets/bucket1/keys/key1
            BODY 	>>	someValue 

*********************************************************************************************
GET OPERATIONS:

/dbs/ <br>
Show All Databases

/dbs/{db}/stats <br>
Returns Grid with Following Information <br>
    Count of Buckets <br>
    Count of Keys <br>
    Size of Database <br>
    Average Bytes Per Key <br>
    (Possibly) Historical Write Speed <br>
    (Possibly) Historical Read Speed <br>

/dbs/{db}/bucket/<br>
Returns list of all buckets<br>
    ** Buckets are automatically created when user attempts to insert a key to a bucket that does not yet exist

/dbs/{db}/bucket/{bucketName}/keys/<br>
Returns list of all keys in bucket

/dbs/{db}/bucket/{bucketName}/keys/{keyName}<br>
Returns {"key":"value"}, returns empty string "" if key does not exist

*********************************************************************************************
DELETE OPERATIONS:
All delete operations require no body and offer no intent confirmation.

/dbs/{db}
Deletes entire database.

/dbs/{db}/bucket/{bucketName}
Deletes bucket & all contents.

/dbs/{db}/bucket/{bucketName}/keys/{keyName}
Deletes key/value pair.

*********************************************************************************************
POST OPERATIONS:

/dbs/{db}/compact
Compacts the database by reading the entire database and writing to a new database.
The original database is then overwritten by the new database. (Reducing size, eliminating deleted keys)