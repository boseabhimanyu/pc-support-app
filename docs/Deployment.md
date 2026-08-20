## Deployment in Vercel

Create a new file at the root of React directory.

`vercel.json`

```
{
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}

```
