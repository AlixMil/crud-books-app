db = db.getSiblingDB('crudbooks');

db.createCollection('books');
db.createCollection('files');
db.createCollection('users');

db.books.createIndex({ title: "text" })
db.books.createIndex({ fileToken: 1 }, { unique: true })
db.files.createIndex({ token: 1 }, { unique: true })
db.users.createIndex({ email: 1 }, { unique: true })

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