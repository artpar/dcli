# dcli

**dcli** is a command-line interface (CLI) tool written in Go for interacting with JSON:API-compliant servers. It allows users to perform CRUD operations, manage relationships, and handle pagination and filtering through a simple command-line interface.


## Prerequisites

- Go (version 1.16 or higher)

## Installation

Clone the repository and navigate to the project directory:

```bash
git clone <repository-url>
cd dcli
```

Build the project:

```bash
go build -o dcli cmd/commands.go
```

This will generate an executable named `dcli`.

## Configuration

Before using the CLI, create a configuration file to specify the base URL of your JSON:API server and other settings.

Create a `config.json` file in your home directory under `.dcli/`:

**Linux/MacOS:**

```bash
mkdir -p ~/.dcli
nano ~/.dcli/config.json
```

**Windows:**

```bash
mkdir %USERPROFILE%\.dcli
notepad %USERPROFILE%\.dcli\config.json
```

Add the following content to `config.json`:

```json
{
  "base_url": "https://your-jsonapi-server.com",
  "api_key": "your-api-key"
}
```

Replace `"https://your-jsonapi-server.com"` with the base URL of your JSON:API server, and provide your API key if required.

## Usage

The `dcli` tool supports various subcommands:

- `create`: Create a new resource.
- `read`: Retrieve a resource by ID.
- `update`: Update an existing resource.
- `delete`: Delete a resource by ID.
- `list`: List resources with optional pagination and filtering.
- `relation`: Manage relationships (get, update, add, remove).

Run `./dcli` without arguments to see the available subcommands.

### Create a Resource

```bash
./dcli create -type=articles -attributes='{"title":"New Article","content":"This is the content."}'
```

- `-type`: The resource type (e.g., `articles`).
- `-attributes`: JSON string of the resource attributes.

### Read a Resource

```bash
./dcli read -type=articles -id=1
```

- `-type`: The resource type.
- `-id`: The ID of the resource.

### Update a Resource

```bash
./dcli update -type=articles -id=1 -attributes='{"title":"Updated Title"}'
```

- `-type`: The resource type.
- `-id`: The ID of the resource.
- `-attributes`: JSON string of the attributes to update.

### Delete a Resource

```bash
./dcli delete -type=articles -id=1
```

- `-type`: The resource type.
- `-id`: The ID of the resource.

### List Resources with Pagination and Filtering

```bash
./dcli list -type=articles -page[number]=1 -page[size]=10 -filter='author:John Doe,category:Tech' -sort='-created_at' -include='comments' -fields='articles:title,content;comments:body'
```

# Updated Documentation

## List API Parameters

The `GET /api/<entityName>` endpoint allows you to retrieve a list of entities of a specified type. You can customize the response using various query parameters to filter, sort, and paginate the results.

### Request Parameters

| Name                 | Type                      | Default Value | Example Value                                                                                               | Description                                                                                                 |
|----------------------|---------------------------|---------------|-------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------|
| `page[number]`       | integer                   | `1`           | `5`                                                                                                         | The page number to retrieve (used for pagination).                                                          |
| `page[size]`         | integer                   | `10`          | `100`                                                                                                       | The number of items per page (used for pagination).                                                         |
| `query`              | JSON Base64               | `[]`          | Base64-encoded JSON array: `[{'column': 'name', 'operator': 'eq', 'value': 'england'}]`                     | Filters the results based on specified conditions.                                                          |
| `group`              | string                    | -             | JSON string: `[{'column': 'name', 'order': 'desc'}]`                                                        | Groups the results based on specified columns.                                                              |
| `included_relations` | comma-separated string    | -             | `user,post,author`                                                                                          | Includes related entities in the response.                                                                  |
| `sort`               | comma-separated string    | -             | `created_at,amount,guest_count`                                                                             | Sorts the results by specified columns. Use `-` prefix for descending order (e.g., `-created_at`).          |
| `filter`             | string                    | -             |                                                                                                             | Filters the results using a simple filter string (implementation-specific).                                 |

#### Parameter Details

- **`page[number]`**: Specifies the page number to retrieve when paginating results. Defaults to `1`.

- **`page[size]`**: Specifies the number of items per page when paginating results. Defaults to `10`.

- **`query`**: A Base64-encoded JSON array of query objects used to filter the results. Each query object can contain:
   - `column`: The column name to filter on.
   - `operator`: The comparison operator (e.g., `eq`, `lt`, `gt`, `like`, etc.).
   - `value`: The value to compare against.

  **Example**:

  To filter entities where the `name` column equals `england`:

  ```json
  [
    {
      "column": "name",
      "operator": "eq",
      "value": "england"
    }
  ]
  ```

  Base64-encode the JSON string and include it in the `query` parameter.

- **`group`**: A JSON string specifying how to group the results.

  **Example**:

  ```json
  [
    {
      "column": "name",
      "order": "desc"
    }
  ]
  ```

- **`included_relations`**: A comma-separated list of related entities to include in the response. For example, `user,post,author`.

- **`sort`**: A comma-separated list of columns to sort the results by. Prefix a column with `-` to sort in descending order.

  **Example**:

   - `sort=created_at` sorts by `created_at` in ascending order.
   - `sort=-created_at,amount` sorts by `created_at` in descending order and then by `amount` in ascending order.

- **`filter`**: A string used to filter results. The specific format and usage depend on the API implementation.

### Example Requests

1. **Basic Pagination**:

   ```http
   GET /api/memory?page[number]=2&page[size]=20
   ```

2. **Filtering with Query**:

   To filter where `name` equals `england`:

   - Prepare the JSON query:

     ```json
     [
       {
         "column": "name",
         "operator": "eq",
         "value": "england"
       }
     ]
     ```

   - Base64-encode the JSON string. Let's say the encoded string is `W3siY29sdW1uIjoibmFtZSIsIm9wZXJhdG9yIjoiZXEiLCJ2YWx1ZSI6ImVuZ2xhbmQifV0=`.

   - Make the request:

     ```http
     GET /api/memory?query=W3siY29sdW1uIjoibmFtZSIsIm9wZXJhdG9yIjoiZXEiLCJ2YWx1ZSI6ImVuZ2xhbmQifV0=
     ```

3. **Including Relations and Sorting**:

   ```http
   GET /api/post?included_relations=author,comments&sort=-created_at
   ```

---

## Describe Command

The `describe` command is used to retrieve and display the schema of a specific entity type, including its columns, relations, and available actions.

### Usage

```bash
./dcli describe -type=<entity_type>
```

### Parameters

- `-type`: Specifies the entity type to describe.

### Example

```bash
./dcli describe -type=user_account
```

### Output

The command will output detailed information about the specified entity, organized into sections:

- **Columns**: Lists the columns (fields) of the entity, including their names, types, data types, and descriptions.
- **Relations**: Lists the relationships the entity has with other entities.
- **Actions**: Lists the actions that can be performed on the entity, along with their input fields.

### Sample Output

```
Entity: user_account

Columns:
Name       Type       Data Type       Description
----       ----       ---------       -----------
password   password   varchar(100)
confirmed  truefalse  bool
name       label      varchar(80)
email      email      varchar(80)
__type     hidden

Relations:
Name                 Relation Type   Related Entity
----                 -------------   --------------
outbox_id            hasMany         outbox
action_id            hasMany         action
artifact_id          hasMany         artifact
...

Actions:
Name                   Label                       Input Fields
----                   -----                       ------------
register_otp           Register Mobile Number      mobile_number(label)    verify_otp      Login with OTP         otp(label), mobile_number(label)
signup                 Sign up                     name(label), email(email), mobile(label), password(password), Password Confirm(password)
...
```

### Notes

- The `describe` command provides a comprehensive overview of an entity's structure and capabilities.


### Manage Relationships

#### Get Relationships

```bash
./dcli relation get -type=articles -id=1 -relation=comments
```

- `-type`: The resource type.
- `-id`: The ID of the resource.
- `-relation`: The name of the relationship.

#### Update Relationships

```bash
./dcli relation update -type=articles -id=1 -relation=comments -data='{"data":[{"type":"comments","id":"2"},{"type":"comments","id":"3"}]}'
```

- `-data`: JSON string representing the relationship data.

#### Add to Relationships

```bash
./dcli relation add -type=articles -id=1 -relation=comments -data='{"data":[{"type":"comments","id":"4"}]}'
```

#### Remove from Relationships

```bash
./dcli relation remove -type=articles -id=1 -relation=comments -data='{"data":[{"type":"comments","id":"2"}]}'
```

## Examples

### Creating an Article

```bash
./dcli create -type=articles -attributes='{"title":"Go CLI Tool","content":"Building a CLI tool in Go."}'
```

### Reading an Article

```bash
./dcli read -type=articles -id=42
```

### Updating an Article

```bash
./dcli update -type=articles -id=42 -attributes='{"title":"Updated Title"}'
```

### Deleting an Article

```bash
./dcli delete -type=articles -id=42
```

### Listing Articles with Filters

```bash
./dcli list -type=articles -filter='author:Jane Doe' -sort='title'
```

### Fetching Article Comments

```bash
./dcli relation get -type=articles -id=42 -relation=comments
```

## Help

For help with a specific command, use the `-h` flag:

```bash
./dcli create -h
./dcli read -h
./dcli update -h
./dcli delete -h
./dcli list -h
./dcli relation -h
./dcli describe -h
./dcli permission -h
```


#### **View Permissions**

```bash
./dcli permission -type=memory -id=123 -action=view
```

**Sample Output:**

```
Permissions for memory (123):
 - GuestPeek
 - UserRead
```

#### **Set Permissions**

```bash
./dcli permission -type=memory -id=123 -action=set -permissions="GuestRead,UserRead,UserUpdate"
```

**Output:**

```
Permissions set successfully.
```

#### **Add Permissions**

```bash
./dcli permission -type=memory -id=123 -action=add -permissions="GuestCreate"
```

**Output:**

```
Permissions added successfully.
```

#### **Remove Permissions**

```bash
./dcli permission -type=memory -id=123 -action=remove -permissions="GuestPeek"
```

**Output:**

```
Permissions removed successfully.
```

---

### **Testing the New Subcommand**

1. **Rebuild the Project**

   ```bash
   go build -o dcli cmd/commands.go
   ```

2. **Test the `view` Action**

   ```bash
   ./dcli permission -type=user_account -id=<object_id> -action=view
   ```

3. **Test the `set` Action**

   ```bash
   ./dcli permission -type=user_account -id=<object_id> -action=set -permissions="UserRead,UserUpdate"
   ```

4. **Test the `add` Action**

   ```bash
   ./dcli permission -type=user_account -id=<object_id> -action=add -permissions="GuestRead"
   ```

5. **Test the `remove` Action**

   ```bash
   ./dcli permission -type=user_account -id=<object_id> -action=remove -permissions="GuestPeek"
   ```



## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For questions or support, please contact [artpar@gmail.com](mailto:artpar@gmail.com).