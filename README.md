
boltapi wraps <a href="https://github.com/boltdb/bolt">boltdb</a> with a REST API interface. 
It allows the user to easily perform CRUD operations for key/value pairs as well as create/delete entire databases.


*********************************************************************************************
PUT OPERATIONS:

Databases & buckets are created (if it does not already exist) when a key/value pair is written using PUT.

/dbs/{db}/bucket/{bucketName}/keys/{keyName} <br>
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
<ul>
<li>Count of Buckets</li>
<li>Count of Keys</li>
<li>Size of Database</li>
<li>Average Bytes Per Key</li>
<li>(Possibly) Historical Write Speed</li>
<li>(Possibly) Historical Read Speed</li>
</ul>

/dbs/{db}/bucket/ <br>
Returns list of all buckets <br>
** Buckets are automatically created when user attempts to insert a key to a bucket that does not yet exist

/dbs/{db}/bucket/{bucketName}/keys/ <br>
Returns list of all keys in bucket

/dbs/{db}/bucket/{bucketName}/keys/{keyName} <br>
Returns {"key":"value"}, returns empty string "" if key does not exist

*********************************************************************************************
DELETE OPERATIONS: <br>
All delete operations require no body and offer no intent confirmation.

/dbs/{db} <br>
Deletes entire database.

/dbs/{db}/bucket/{bucketName} <br>
Deletes bucket & all contents.

/dbs/{db}/bucket/{bucketName}/keys/{keyName} <br>
Deletes key/value pair.

*********************************************************************************************
POST OPERATIONS:

/dbs/{db}/compact <br>
Compacts the database by reading the entire database and writing to a new database. <br>
The original database is then overwritten by the new database. (Reducing size, eliminating deleted keys)