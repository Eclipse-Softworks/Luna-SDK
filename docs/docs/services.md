# Service Modules

## 1. Core Services (Users & Projects)

Manage the foundational entities of the platform.

### TypeScript
```typescript
// List users
const users = await client.users.iterate({ limit: 50 });
for await (const user of users) {
    console.log(user.email);
}

// Create project
const proj = await client.projects.create({ name: 'Alpha', description: 'Main project' });
```

## 2. ResMate (Student Residences)

Search and filter student accommodations.

### Python
```python
# Search for NSFAS accredited residences near a campus
residences = await client.res_mate.residences.iterate(
    nsfas=True,
    campus_id="cmp_uct_main",
    max_price=6000
)

async for res in residences:
    print(f"{res.name} - {res.rating} stars")
```

## 3. ES Identity (Groups & Permissions)

Manage access control groups.

### Go
```go
// Create a new group
group, err := client.Identity().Groups().Create(ctx, &luna.GroupCreate{
    Name: "Admins",
    Permissions: []string{"read", "write", "delete"},
})
```

## 4. ES Storage (Buckets & Files)

Manage cloud storage buckets and file uploads.

### TypeScript
```typescript
// Upload a file (Note: In a browser, body would be a File/Blob)
const file = await client.storage.buckets.upload('bkt_images', {
    name: 'profile.jpg',
    content: blobData
});

// Get download URL
const url = await client.storage.files.getDownloadUrl(file.id);
```

## 5. ES AI Tools

Access Generative AI models.

### Python
```python
response = await client.ai.chat_completions(
    model="gpt-4",
    messages=[
        {"role": "system", "content": "You are a coding assistant."},
        {"role": "user", "content": "Explain async/await."}
    ]
)
print(response.choices[0].message.content)
```

## 6. ES Automation

Trigger and manage workflows.

### Go
```go
// Trigger a deployment workflow
run, err := client.Automation().Workflows().Trigger(ctx, "wf_deploy", map[string]interface{}{
    "environment": "production",
    "version": "v1.2.3",
})
fmt.Printf("Workflow started: %s\n", run.ID)
```
