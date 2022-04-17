# List all users

## Definition
`GET /api/v1/users`

## Response
- 200 OK on success
```json
[
    {
        "id": 1,
        "name": "Maria"
    },
    {
        "id": 2,
        "name": "John"
    }
]
```

# Register a new user

## Definition
`PUT /api/v1/users`

## Response
- 201 CREATED on success
```json
{
    "name": "Joe"
}
```

# List a user

## Definition
`GET /api/v1/users/<id>`

## Response
- 200 OK on success
- 404 NOT FOUND if the user with that id does not exist 
```json
{
    "id": 1,
    "name": "Maria"
}
```

# Habits
## Properties
```
[
    {
        "id": 1
        "name": "Read 📚",
        "days": ["Saturday", "Sunday"] 
    },
    {
        "id": 2,
        "name": "Play guitar 🎸"
        "days": ["Monday", "Sunday"] 
    }
]
```