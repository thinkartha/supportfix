# Fix 403 Forbidden and 404 on API

## 404 "page not found" on raw API URL

If `POST https://xxx.execute-api.us-east-1.amazonaws.com/prod/api/auth/login` returns **404**, API Gateway includes the stage (`prod`) in the path, but our routes expect `/api/...`. The Lambda handler strips the stage prefix before routing. **Redeploy the backend** so the updated Lambda is deployed.

## 403 Forbidden on https://api.supportfix.ai

When `POST https://api.supportfix.ai/api/auth/login` returns `{"message": "Forbidden"}`, the response is from **API Gateway**, not your Lambda. Common causes and fixes:

---

## 1. Check the raw API Gateway URL first

Test the Lambda directly to confirm it works:

```bash
# Get your API invoke URL from CloudFormation
aws cloudformation describe-stacks --stack-name supportdesk-api \
  --query "Stacks[0].Outputs[?OutputKey=='ApiUrl'].OutputValue" --output text
```

Then in Postman, change **baseUrl** to that URL (e.g. `https://abc123.execute-api.us-east-1.amazonaws.com/prod`) and try **Login** again.

- **If it works:** The issue is with the custom domain setup (see step 2).
- **If it fails:** The issue is with the Lambda/API itself.

---

## 2. Fix custom domain API mapping

403 often means the custom domain has **no API mapping** or the mapping is wrong.

1. **API Gateway** → **Custom domain names**
2. Select **api.supportfix.ai**
3. Open **API mappings**
4. Ensure there is a mapping:
   - **API:** Your SAM HTTP API (e.g. `supportdesk-api` or the API ID)
   - **Stage:** `prod`
   - **Path:** Leave blank (or `$default`) so `api.supportfix.ai/api/...` goes to your API

5. If there is no mapping, click **Configure API mappings** → **Add new mapping** and add the above.

---

## 3. Verify path (no stage in URL)

When using the custom domain, do **not** include the stage in the path:

| Correct                          | Incorrect                           |
|----------------------------------|-------------------------------------|
| `https://api.supportfix.ai/api/auth/login` | `https://api.supportfix.ai/prod/api/auth/login` |

---

## 4. Check Route 53

- **Route 53** → **Hosted zones** → **supportfix.ai**
- Confirm there is a record for **api** pointing to the API Gateway custom domain target.
- See `docs/ROUTE53_API_SUPPORTFIX.md` for setup.

---

## 5. ACM certificate

The custom domain needs an ACM certificate in **us-east-1** that covers `api.supportfix.ai` or `*.supportfix.ai`. Certificate status should be **Issued**.
