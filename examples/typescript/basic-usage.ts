/**
 * Luna SDK - TypeScript Basic Usage Examples
 *
 * This file demonstrates common operations with the Luna SDK.
 * Run with: npx tsx basic-usage.ts
 */

import { LunaClient } from '@eclipse-softworks/luna-sdk';

// Initialize the client with API key authentication
const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY!,
    // Optional configuration
    timeout: 30000,
    maxRetries: 3,
});

// ============================================
// Example 1: User Management
// ============================================

async function userManagementExample() {
    console.log('=== User Management ===\n');

    // List users with pagination
    const userList = await client.users.list({ limit: 10 });
    console.log(`Found ${userList.data.length} users`);

    // Create a new user
    const newUser = await client.users.create({
        email: 'jane.doe@example.com',
        name: 'Jane Doe',
    });
    console.log(`Created user: ${newUser.id}`);

    // Get user details
    const user = await client.users.get(newUser.id);
    console.log(`User name: ${user.name}, Email: ${user.email}`);

    // Update the user
    const updatedUser = await client.users.update(newUser.id, {
        name: 'Jane M. Doe',
        avatar_url: 'https://example.com/avatar.jpg',
    });
    console.log(`Updated user name: ${updatedUser.name}`);

    // Delete the user
    await client.users.delete(newUser.id);
    console.log('User deleted');
}

// ============================================
// Example 2: Project Management
// ============================================

async function projectManagementExample() {
    console.log('\n=== Project Management ===\n');

    // Create a project
    const project = await client.projects.create({
        name: 'My Awesome App',
        description: 'A revolutionary application built with Luna SDK',
    });
    console.log(`Created project: ${project.id}`);

    // List all projects
    const projects = await client.projects.list({ limit: 20 });
    console.log(`Total projects: ${projects.data.length}`);

    // Get project details
    const projectDetails = await client.projects.get(project.id);
    console.log(`Project: ${projectDetails.name}`);
    console.log(`Owner: ${projectDetails.owner_id}`);
    console.log(`Created: ${projectDetails.created_at}`);

    // Update project
    const updated = await client.projects.update(project.id, {
        description: 'Updated description with new features',
    });
    console.log(`Updated project: ${updated.description}`);

    // Clean up
    await client.projects.delete(project.id);
    console.log('Project deleted');
}

// ============================================
// Example 3: Paginating Through Results
// ============================================

async function paginationExample() {
    console.log('\n=== Pagination Example ===\n');

    // Using the iterator for automatic pagination
    let count = 0;
    for await (const user of client.users.iterate({ limit: 10 })) {
        console.log(`User: ${user.name} (${user.email})`);
        count++;
        if (count >= 50) break; // Limit for demo
    }
    console.log(`Iterated through ${count} users`);
}

// ============================================
// Example 4: Error Handling
// ============================================

async function errorHandlingExample() {
    console.log('\n=== Error Handling ===\n');

    try {
        // Try to get a non-existent user
        await client.users.get('usr_nonexistent123');
    } catch (error) {
        if (error instanceof Error) {
            console.log(`Error caught: ${error.message}`);
            // Check for specific error types
            if ('code' in error) {
                console.log(`Error code: ${(error as any).code}`);
            }
        }
    }
}

// ============================================
// Run all examples
// ============================================

async function main() {
    try {
        await userManagementExample();
        await projectManagementExample();
        await paginationExample();
        await errorHandlingExample();
    } catch (error) {
        console.error('Example failed:', error);
        process.exit(1);
    }
}

main();
