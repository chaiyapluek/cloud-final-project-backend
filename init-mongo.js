db.createUser(
    {
        user: "admin",
        pwd: "123456",
        roles: [
            {
                role: "readWrite",
                db: "saywub"
            }
        ]
    }
);
db.createCollection("login_attempts");
db.createCollection("users");
db.createCollection("locations");
db.createCollection("menus");

db.locations.insertMany([
    {
        "_id": ObjectId("651986655ff54c62ea7ff99c"),
        "name": "วิศวจุฬา",
    }
]);

db.menus.insertMany([
    {
        "location_id": ObjectId("651986655ff54c62ea7ff99c"),
        "name": "[Promotion] Sandwich",
        "description": "Today only! Get 100% off on all sandwiches",
        "price": 129,
        "iconImage": "/static/image/classic.jpg",
        "thumbnailImage": "/static/image/classic-thumbnail.jpg",
        "steps": [
            {
                "name": "Bread",
                "description": "Choose your bread",
                "type": "radio",
                "required": true,
                "min": 1,
                "max": 1,
                "options": [
                    {
                        "name": "White",
                        "value": "white",
                        "price": 0
                    },
                    {
                        "name": "Brown",
                        "value": "brown",
                        "price": 0
                    }
                ]
            },
            {
                "name": "Protein",
                "description": "Choose your protein, (choose 1 to 2)",
                "type": "checkbox",
                "required": true,
                "min": 1,
                "max": 2,
                "options": [
                    {
                        "name": "Chicken",
                        "value": "chicken",
                        "price": 0
                    },
                    {
                        "name": "Beef",
                        "value": "beef",
                        "price": 39
                    },
                    {
                        "name": "Tuna",
                        "value": "tuna",
                        "price": 29
                    }
                ]
            },
            {
                "name": "Vegetables",
                "description": "Choose your vegetables, (choose 0 to 3)",
                "type": "checkbox",
                "required": false,
                "min": 0,
                "max": 3,
                "options": [
                    {
                        "name": "Lettuce",
                        "value": "lettuce",
                        "price": 0
                    },
                    {
                        "name": "Tomato",
                        "value": "tomato",
                        "price": 0
                    },
                    {
                        "name": "Onion",
                        "value": "onion",
                        "price": 0
                    }
                ]
            },
            {
                "name": "Sauces",
                "description": "Choose your sauces, (choose 0 to 2)",
                "type": "checkbox",
                "required": false,
                "min": 0,
                "max": 2,
                "options": [
                    {
                        "name": "Mayo",
                        "value": "mayo",
                        "price": 0
                    },
                    {
                        "name": "Ketchup",
                        "value": "ketchup",
                        "price": 0
                    },
                    {
                        "name": "Mustard",
                        "value": "mustard",
                        "price": 0
                    }
                ]
            }
        ]
    }
]);