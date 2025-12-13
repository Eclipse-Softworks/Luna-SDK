/**
 * Luna SDK - File Storage Example
 *
 * Build a file management system using Luna's Storage API.
 * This example demonstrates bucket management, file uploads, and downloads.
 */

import { LunaClient } from '@eclipse-softworks/luna-sdk';
import * as fs from 'fs';
import * as path from 'path';

const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY!,
});

// ============================================
// Storage Manager Class
// ============================================

class StorageManager {
    private defaultBucketId: string | null = null;

    /**
     * List all available storage buckets
     */
    async listBuckets() {
        console.log('Listing storage buckets...\n');

        const buckets = await client.storage.buckets.list();

        console.log(`Found ${buckets.data.length} buckets:\n`);
        buckets.data.forEach((bucket) => {
            console.log(`  • ${bucket.name} (${bucket.id})`);
            console.log(`    Region: ${bucket.region}`);
        });

        return buckets.data;
    }

    /**
     * Get bucket details
     */
    async getBucket(bucketId: string) {
        const bucket = await client.storage.buckets.get(bucketId);
        console.log(`\nBucket: ${bucket.name}`);
        console.log(`   ID: ${bucket.id}`);
        console.log(`   Region: ${bucket.region}`);
        return bucket;
    }

    /**
     * List files in a bucket
     */
    async listFiles(bucketId: string) {
        console.log(`\nFiles in bucket ${bucketId}:\n`);

        const files = await client.storage.files.list(bucketId);

        if (files.data.length === 0) {
            console.log('   (No files found)');
            return [];
        }

        files.data.forEach((file) => {
            const sizeKB = (file.size / 1024).toFixed(2);
            console.log(`  [FILE] ${file.key}`);
            console.log(`     ID: ${file.id}`);
            console.log(`     Size: ${sizeKB} KB`);
            console.log(`     Type: ${file.content_type}`);
        });

        return files.data;
    }

    /**
     * Upload a file to a bucket
     */
    async uploadFile(
        bucketId: string,
        filePath: string,
        metadata?: Record<string, string>
    ) {
        console.log(`\nUploading ${filePath}...`);

        // Read file content
        const fileContent = fs.readFileSync(filePath);
        const fileName = path.basename(filePath);

        // Upload to Luna Storage
        const uploadedFile = await client.storage.buckets.upload(
            bucketId,
            fileContent,
            {
                filename: fileName,
                ...metadata,
            }
        );

        console.log(`File uploaded successfully!`);
        console.log(`   ID: ${uploadedFile.id}`);
        console.log(`   Key: ${uploadedFile.key}`);
        console.log(`   URL: ${uploadedFile.url}`);

        return uploadedFile;
    }

    /**
     * Get download URL for a file
     */
    async getDownloadUrl(fileId: string) {
        const result = await client.storage.files.getDownloadUrl(fileId);
        console.log(`\nDownload URL: ${result.url}`);
        return result.url;
    }

    /**
     * Download a file
     */
    async downloadFile(fileId: string, outputPath: string) {
        console.log(`\nDownloading file ${fileId}...`);

        const { url } = await client.storage.files.getDownloadUrl(fileId);

        // Fetch the file content
        const response = await fetch(url);
        const buffer = await response.arrayBuffer();

        // Save to disk
        fs.writeFileSync(outputPath, Buffer.from(buffer));

        console.log(`File downloaded to: ${outputPath}`);
    }
}

// ============================================
// Example: Document Management System
// ============================================

async function documentManagementExample() {
    console.log('Document Management System\n');
    console.log('='.repeat(50));

    const storage = new StorageManager();

    // List available buckets
    const buckets = await storage.listBuckets();

    if (buckets.length === 0) {
        console.log('No buckets available. Please create a bucket first.');
        return;
    }

    // Use the first bucket
    const bucket = buckets[0];
    console.log(`\nUsing bucket: ${bucket.name}`);

    // List existing files
    await storage.listFiles(bucket.id);
}

// ============================================
// Example: Profile Picture Upload
// ============================================

async function profilePictureExample() {
    console.log('\nProfile Picture Upload Example\n');
    console.log('='.repeat(50));

    const storage = new StorageManager();

    // Simulated upload (in real app, you'd have actual file)
    const buckets = await storage.listBuckets();

    if (buckets.length > 0) {
        const bucket = buckets[0];

        // Example: Upload a profile picture
        // In a real application:
        // const uploaded = await storage.uploadFile(
        //     bucket.id,
        //     './profile.jpg',
        //     { userId: 'usr_123', type: 'avatar' }
        // );

        console.log(`
To upload a file, use:

const uploaded = await client.storage.buckets.upload(
    '${bucket.id}',
    fileContent,
    { filename: 'avatar.jpg' }
);

console.log('Uploaded:', uploaded.url);
`);
    }
}

// ============================================
// Example: Batch File Operations
// ============================================

async function batchOperationsExample() {
    console.log('\nBatch File Operations Example\n');
    console.log('='.repeat(50));

    const storage = new StorageManager();
    const buckets = await storage.listBuckets();

    if (buckets.length === 0) {
        console.log('No buckets available.');
        return;
    }

    const bucket = buckets[0];

    // List all files and generate download URLs
    const files = await storage.listFiles(bucket.id);

    if (files.length > 0) {
        console.log('\nGenerating download URLs for all files:\n');

        for (const file of files) {
            const url = await storage.getDownloadUrl(file.id);
            console.log(`  ${file.key}: ${url}`);
        }
    }
}

// ============================================
// Example: Storage Usage Report
// ============================================

async function storageReportExample() {
    console.log('\nStorage Usage Report\n');
    console.log('='.repeat(50));

    const storage = new StorageManager();
    const buckets = await storage.listBuckets();

    let totalFiles = 0;
    let totalSize = 0;

    for (const bucket of buckets) {
        console.log(`\n${bucket.name} (${bucket.region})`);

        const files = await client.storage.files.list(bucket.id);
        const bucketSize = files.data.reduce((sum, file) => sum + file.size, 0);

        totalFiles += files.data.length;
        totalSize += bucketSize;

        console.log(`   Files: ${files.data.length}`);
        console.log(`   Size: ${(bucketSize / 1024 / 1024).toFixed(2)} MB`);
    }

    console.log('\n' + '='.repeat(50));
    console.log(`Total Storage Used`);
    console.log(`   Buckets: ${buckets.length}`);
    console.log(`   Files: ${totalFiles}`);
    console.log(`   Total Size: ${(totalSize / 1024 / 1024).toFixed(2)} MB`);
}

// ============================================
// Main
// ============================================

async function main() {
    try {
        await documentManagementExample();
        await profilePictureExample();
        await batchOperationsExample();
        await storageReportExample();
    } catch (error) {
        console.error('Error:', error);
        process.exit(1);
    }
}

main();
