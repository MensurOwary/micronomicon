db = db.getSiblingDB('admin')

db.auth('root', 'root')

db = db.getSiblingDB('micron')

db.jwts.createIndex( {createdDate: 1}, {
    expireAfterSeconds: 2 * 3600 // 2 hours
});