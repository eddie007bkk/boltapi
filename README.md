
        *********************************************************************************************
		PUT OPERATIONS:
		Databases & buckets are created (if it does not already exist) when a key/value pair is written using PUT.

		/dbs/{db}/bucket/{bucketName}/keys/{keyName}
			PUT - Creates Database and/or bucket if they do not already exist. Inserts key/value pair using BODY as value.

			Ex:
				 	PUT 	>> 	http://host:8080/dbs/somenewdatabase.db/buckets/bucket1/keys/key1
					BODY 	>>	someValue 

		*********************************************************************************************
		GET OPERATIONS:

		/dbs/
			GET - Show All Databases

		/dbs/{db}/stats
			GET - Returns Grid with Following Information
					Count of Buckets
					Count of Keys
					Size of Database
					Average Bytes Per Key
					(Possibly) Historical Write Speed
					(Possibly) Historical Read Speed

		/dbs/{db}/bucket/
			GET - Returns list of all buckets
			** Buckets are automatically created when user attempts to insert a key to a bucket that does not yet exist

		/dbs/{db}/bucket/{bucketName}/keys/
			GET - Returns list of all keys in bucket

		/dbs/{db}/bucket/{bucketName}/keys/{keyName}
			GET - Returns {"key":"value"}, returns empty string "" if key does not exist
		
		*********************************************************************************************
		DELETE OPERATIONS:
		All delete operations require no body and offer no intent confirmation.

		/dbs/{db}
			DELETE - Deletes entire database.

		/dbs/{db}/bucket/{bucketName}
			DELETE - Deletes bucket & all contents.

		/dbs/{db}/bucket/{bucketName}/keys/{keyName}
			DELETE - Deletes key/value pair.

		*********************************************************************************************
		POST OPERATIONS:

		/dbs/{db}/compact
			POST - Compacts the database by reading the entire database and writing to a new database.
				  The original database is then overwritten by the new database. (Reducing size, eliminating deleted keys)