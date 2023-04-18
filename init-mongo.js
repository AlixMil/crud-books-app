db = db.getSiblingDB('crudbooks');

db.createCollection('books');
db.createCollection('files');
db.createCollection('users');

db.createUser(
	{
		user: "admin",
		pwd: "0000",
		roles: [
			{
				role: "readWrite",
				db: "crudbooks"
			}
		]
	}
);