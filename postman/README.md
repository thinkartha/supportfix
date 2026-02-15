# SupportDesk API – Postman collection

Use this collection to test the SupportDesk API via API Gateway.

## Import

1. Open Postman.
2. **Import** → **Upload Files** → select `SupportDesk-API.postman_collection.json`.

## Setup

1. Open the collection **SupportDesk API (API Gateway)**.
2. Go to the **Variables** tab.
3. Set **baseUrl** to your API Gateway URL, e.g.  
   `https://YOUR_API_ID.execute-api.us-east-1.amazonaws.com/prod`  
   (Get this from the Lambda deploy output or AWS Console → API Gateway → your API → Stages → prod → Invoke URL.)
4. Leave **token** empty; it will be set automatically after **Login**.

## Usage

1. Run **Auth → Login** with a valid `email` and `password`.  
   The response includes a JWT; the collection script saves it into **token**.
2. All other folders (Users, Organizations, Tickets, etc.) use **Authorization: Bearer {{token}}**.
3. For requests with **:id** in the path (e.g. Get user by ID), set the **Params** tab **path** variable `id` to the real ID (e.g. `user-xxx`, `ticket-xxx`).

## Folders

| Folder         | Endpoints |
|----------------|-----------|
| Auth           | Login, Me (current user) |
| Health         | Health check |
| Users          | List, Get, Create, Update, Delete |
| Organizations  | List, Get, Create, Update, Delete |
| Tickets        | List, Get, Create, Update, Add message, Add time entry, Request conversion |
| Approvals      | List, Update |
| Invoices       | List, Create, Update status |
| Dashboard      | Stats, Activities |
