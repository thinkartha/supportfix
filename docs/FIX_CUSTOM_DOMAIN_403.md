# Fix 403 Forbidden on https://api.supportfix.ai

The raw URL works (`https://340su058u8.execute-api.us-east-1.amazonaws.com/prod/api/auth/login`) but the custom domain returns 403. This means the **API mapping** is missing or wrong.

---

## 1. Check current setup

Run these commands (replace region if needed):

```bash
# List custom domains
aws apigatewayv2 get-domain-names --region us-east-1

# Get your API ID from the working URL (it's 340su058u8)
API_ID=340su058u8

# List API mappings for api.supportfix.ai
aws apigatewayv2 get-api-mappings --domain-name api.supportfix.ai --region us-east-1
```

If `get-api-mappings` returns empty `Items: []` or errors, the mapping is missing.

---

## 2. Create the API mapping

```bash
# Create API mapping: api.supportfix.ai → your HTTP API, prod stage
aws apigatewayv2 create-api-mapping \
  --domain-name api.supportfix.ai \
  --api-id 340su058u8 \
  --stage prod \
  --region us-east-1
```

**Leave `--api-mapping-key` out** so the root path (`/`) maps to your API. Then:
- `https://api.supportfix.ai/api/auth/login` → your Lambda

---

## 3. If the custom domain doesn't exist yet

Create it first:

```bash
# 1. Get ACM certificate ARN for api.supportfix.ai (must be in us-east-1)
aws acm list-certificates --region us-east-1 --query "CertificateSummaryList[?DomainName=='*.supportfix.ai' || DomainName=='api.supportfix.ai'].CertificateArn" --output text

# 2. Create custom domain (use the ARN from step 1)
aws apigatewayv2 create-domain-name \
  --domain-name api.supportfix.ai \
  --domain-name-configurations "CertificateArn=arn:aws:acm:us-east-1:YOUR_ACCOUNT:certificate/YOUR_CERT_ID" \
  --region us-east-1

# 3. Add API mapping (from step 2 above)
```

---

## 4. Verify Route 53

Ensure `api.supportfix.ai` points to the API Gateway custom domain target:

```bash
# Get the target for your custom domain
aws apigatewayv2 get-domain-name --domain-name api.supportfix.ai --region us-east-1 \
  --query "DomainNameConfigurations[0].ApiGatewayDomainName" --output text
```

Route 53 record `api` in hosted zone `supportfix.ai` should be a **CNAME** or **A (Alias)** to that target (e.g. `d-xxxxxx.execute-api.us-east-1.amazonaws.com`).

---

## 5. Test after fix

```bash
curl -s -o /dev/null -w "%{http_code}" https://api.supportfix.ai/api/health
# Should return 200

curl -X POST https://api.supportfix.ai/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@supportfix.ai","password":"admin123"}'
# Should return 200 with token
```
