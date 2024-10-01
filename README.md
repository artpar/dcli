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

- `-type`: The resource type.
- `-page[number]`: (Optional) Page number.
- `-page[size]`: (Optional) Number of items per page.
- `-filter`: (Optional) Filters in `key1:value1,key2:value2` format.
- `-sort`: (Optional) Fields to sort by (e.g., `name,-created_at`).
- `-include`: (Optional) Related resources to include.
- `-fields`: (Optional) Sparse fieldsets in `type1:field1,field2;type2:field3` format.

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
```

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For questions or support, please contact [your-email@example.com](mailto:your-email@example.com).