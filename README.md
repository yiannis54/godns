# Go-DNS

**Dynamically update a domain to target your local server through Cloudflare**


## Prerequisites

**A Type DNS record in cloudflare is required**

- Login or create an account on [Cloudflare](https://www.cloudflare.com/)
- Add your domain and create a simple DNS `A Type` (No proxied) record pointing to your IP. https://developers.cloudflare.com/dns/manage-dns-records/how-to/create-zone-apex/
![2024-11-03_11-49](https://github.com/user-attachments/assets/97648c8e-bac1-42f6-850e-ce4ea73722c7)

**An API Token is required**

- Go to [https://dash.cloudflare.com/profile/api-tokens](https://dash.cloudflare.com/profile/api-tokens) and create a new Token
    - Token should have permissions: Zone = DNS (edit)
    - Resources `Include = YOUR_DOMAIN`

## Usage

- Create a `.env` file based on the example in `.env.dist`.
    - Insert your `Domain` and `Cloudflare API Token` in the corresponding environment variables.
    - Update interval is `1 hour` and server port is `3003`, change to your own needs.
    - Default login password is `admin`.
    - If Sentry DNS is filled, it will track any errors logged from the process.
- Run the application with `export $(cat .env) && make run`

You can make an executable for a Linux machine with `make build` and use it to create a service at your server.

## View DNS record changes

You can visit `localhost:3003/records/` to check on IP changes. Default username/password is admin/admin. You can change the password to your own needs in the environment variables.
![2024-11-03_11-43](https://github.com/user-attachments/assets/894d5b8e-6a25-4598-8deb-0acd5b8bbd4e)
