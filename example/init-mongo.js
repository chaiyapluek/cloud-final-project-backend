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
        "price": 89,
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
                        "name": "Wheat",
                        "value": "wheat",
                        "price": 0
                    },
                    {
                        "name": "Honey Oat",
                        "value": "honey_oat",
                        "price": 0
                    },
                    {
                        "name": "Italian",
                        "value": "italian",
                        "price": 0
                    },
                    {
                        "name": "Parmesan Oregano",
                        "value": "parmesan_oregano",
                        "price": 0
                    },
                    {
                        "name": "Flatbread",
                        "value": "flatbread",
                        "price": 0
                    }
                ]
            },
            {
                "name": "Size",
                "description": "Choose your sandwich size",
                "type": "radio",
                "required": true,
                "min": 1,
                "max": 1,
                "options": [
                    {
                        "name": "6 inch",
                        "value": "6_inch",
                        "price": 0
                    },
                    {
                        "name": "Footlong",
                        "value": "footlong",
                        "price": 59
                    }
                ]
            },
            {
                "name": "Menu",
                "description": "Choose you favorite menu",
                "type": "radio",
                "required": true,
                "min": 1,
                "max": 1,
                "options": [
                    {
                        "name": "Teriyaki Chicken",
                        "value": "teriyaki_chicken",
                        "price": 49
                    },
                    {
                        "name": "Roasted Beef",
                        "value": "roasted_beef",
                        "price": 49
                    },
                    {
                        "name": "BBQ Chicken",
                        "value": "bbq_chicken",
                        "price": 39
                    },
                    {
                        "name": "Roasted Chicken",
                        "value": "roasted_chicken",
                        "price": 29
                    },
                    {
                        "name": "Veggie Delux",
                        "value": "veggie_delux",
                        "price": 29
                    },
                    {
                        "name": "Tuna",
                        "value": "tuna",
                        "price": 10
                    },
                    {
                        "name": "Slice Chicken",
                        "value": "slice_chicken",
                        "price": 10
                    },
                    {
                        "name": "Veggie Delight",
                        "value": "veggie_delight",
                        "price": 0
                    }
                ]
            },
            {
                "name": "Vegetables",
                "description": "Choose your vegetables, (choose 1 or more)",
                "type": "checkbox",
                "required": true,
                "min": 1,
                "max": 99,
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
                        "name": "Cumcumber",
                        "value": "cumcumber",
                        "price": 0
                    },
                    {
                        "name": "Pickles",
                        "value": "pickles",
                        "price": 0
                    },
                    {
                        "name": "Green Peppers",
                        "value": "green_peppers",
                        "price": 0
                    },
                    {
                        "name": "Olives",
                        "value": "olives",
                        "price": 0
                    },
                    {   
                        "name": "Onion",
                        "value": "onion",
                        "price": 0
                    },
                    {
                        "name": "Jalapenos",
                        "value": "jalapenos",
                        "price": 0
                    }
                ]
            },
            {
                "name": "Sauces",
                "description": "Choose your sauces, (choose between 0 to 3)",
                "type": "checkbox",
                "required": true,
                "min": 0,
                "max": 3,
                "options": [
                    {
                        "name": "Honey Mustard",
                        "value": "honey_mustard",
                        "price": 0
                    },
                    {
                        "name": "Sweet Onion",
                        "value": "sweet_onion",
                        "price": 0
                    },
                    {
                        "name": "Chipotle Southwest",
                        "value": "chipotle_southwest",
                        "price": 0
                    },
                    {
                        "name": "Mayonnaise",
                        "value": "mayonnaise",
                        "price": 0
                    },
                    {
                        "name": "BBQ Sauce",
                        "value": "bbq_sauce",
                        "price": 0
                    },
                    {
                        "name": "Tomato Sauce",
                        "value": "tomato_sauce",
                        "price": 0
                    },
                    {
                        "name": "Thousand Island Dressing",
                        "value": "thousand_island_dressing",
                        "price": 0
                    },
                    {
                        "name": "Hot chili sauce",
                        "value": "hot_chili_sauce",
                        "price": 0
                    }
                ]
            },
            {
                "name": "Add-ons",
                "description": "Choose your add-ons, (choose 0 or more)",
                "type": "checkbox",
                "required": false,
                "min": 0,
                "max": 99,
                "options": [
                    {
                        "name": "Mozzarella Cheese",
                        "value": "mozzarella_cheese",
                        "price": 15
                    },
                    {
                        "name": "Cheddar Cheese",
                        "value": "cheddar_cheese",
                        "price": 15
                    },
                    {
                        "name": "Bacon",
                        "value": "bacon",
                        "price": 40
                    },
                    {
                        "name": "Egg",
                        "value": "egg",
                        "price": 15
                    },
                    {
                        "name": "Double Meat",
                        "value": "double_mear",
                        "price": 60
                    },
                    {
                        "name": "Avocado",
                        "value": "avocado",
                        "price": 40
                    },
                    {
                        "name": "Chopped Mushroom",
                        "value": "chopped_mushroom",
                        "price": 20
                    }
                ]
            },
            {
                "name": "Meal",
                "description": "Make it a meal!",
                "type": "checkbox",
                "required": false,
                "min": 0,
                "max": 1,
                "options": [
                    {
                        "name": "Cookie + 22 Oz. Pepsi",
                        "value": "cookie_pepsi",
                        "price": 59
                    }
                ]
            }
        ]
    }
]);
