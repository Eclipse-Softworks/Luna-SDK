/**
 * Luna SDK - Workflow Automation Example
 *
 * Build automated workflows using Luna's Automation API.
 * This example demonstrates workflow management and triggering.
 */

import { LunaClient } from '@eclipse-softworks/luna-sdk';

const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY!,
});

// ============================================
// Workflow Manager Class
// ============================================

class WorkflowManager {
    /**
     * List all available workflows
     */
    async listWorkflows() {
        console.log('Available Workflows\n');

        const workflows = await client.automation.workflows.list();

        workflows.data.forEach((workflow) => {
            const status = workflow.is_active ? '[Active]' : '[Inactive]';
            console.log(`${status} ${workflow.name}`);
            console.log(`   ID: ${workflow.id}`);
            console.log(`   Trigger: ${workflow.trigger_type}`);
            console.log('');
        });

        return workflows.data;
    }

    /**
     * Trigger a workflow
     */
    async triggerWorkflow(workflowId: string, params: Record<string, unknown> = {}) {
        console.log(`\nTriggering workflow ${workflowId}...`);

        const run = await client.automation.workflows.trigger(workflowId, params);

        console.log(`Workflow triggered!`);
        console.log(`   Run ID: ${run.id}`);
        console.log(`   Status: ${run.status}`);
        console.log(`   Started: ${run.started_at}`);

        return run;
    }

    /**
     * Get active workflows
     */
    async getActiveWorkflows() {
        const workflows = await client.automation.workflows.list();
        return workflows.data.filter((w) => w.is_active);
    }
}

// ============================================
// Example: User Onboarding Workflow
// ============================================

async function userOnboardingExample() {
    console.log('User Onboarding Automation\n');
    console.log('='.repeat(50));

    const manager = new WorkflowManager();

    // List available workflows
    await manager.listWorkflows();

    // Example: Trigger onboarding workflow for new user
    console.log('\nTriggering user onboarding workflow...\n');

    // In a real application, you'd trigger with actual user data
    const mockWorkflowId = 'wf_onboarding_123';
    const newUserData = {
        userId: 'usr_new_456',
        email: 'newuser@example.com',
        name: 'New User',
        plan: 'premium',
    };

    console.log('Workflow parameters:', newUserData);
    console.log(`
// Trigger the workflow:
const run = await client.automation.workflows.trigger(
    '${mockWorkflowId}',
    ${JSON.stringify(newUserData, null, 4)}
);
`);
}

// ============================================
// Example: Scheduled Report Generation
// ============================================

async function reportGenerationExample() {
    console.log('\nReport Generation Workflow\n');
    console.log('='.repeat(50));

    const manager = new WorkflowManager();

    // Get workflows that are scheduled
    const workflows = await manager.listWorkflows();
    const scheduledWorkflows = workflows.filter((w) => w.trigger_type === 'schedule');

    console.log(`Found ${scheduledWorkflows.length} scheduled workflows`);

    // Example trigger for report generation
    console.log(`
// Manually trigger a scheduled report:
const run = await client.automation.workflows.trigger('wf_weekly_report', {
    reportType: 'usage',
    dateRange: {
        start: '2024-01-01',
        end: '2024-01-31'
    },
    recipients: ['team@example.com'],
    format: 'pdf'
});
`);
}

// ============================================
// Example: Event-Driven Workflows
// ============================================

async function eventDrivenExample() {
    console.log('\nEvent-Driven Workflows\n');
    console.log('='.repeat(50));

    console.log(`
Event-driven workflows are automatically triggered when specific events occur.
Examples of events that can trigger workflows:

- user.created - New user registration
- user.deleted - User account deletion
- project.created - New project creation
- file.uploaded - File upload to storage
- payment.completed - Payment processed

Example: Setting up an event-driven notification:

// When a user is created, this workflow sends a welcome email
const workflows = await client.automation.workflows.list();
const welcomeWorkflow = workflows.data.find(
    w => w.name === 'Welcome Email' && w.trigger_type === 'event'
);

// Event workflows are triggered automatically, but can be manually triggered:
if (welcomeWorkflow) {
    await client.automation.workflows.trigger(welcomeWorkflow.id, {
        event: 'user.created',
        data: {
            userId: 'usr_123',
            email: 'user@example.com'
        }
    });
}
`);
}

// ============================================
// Example: Webhook Integration
// ============================================

async function webhookExample() {
    console.log('\nWebhook Workflows\n');
    console.log('='.repeat(50));

    const manager = new WorkflowManager();
    const workflows = await manager.listWorkflows();

    const webhookWorkflows = workflows.filter((w) => w.trigger_type === 'webhook');
    console.log(`Found ${webhookWorkflows.length} webhook-triggered workflows`);

    console.log(`
Webhook workflows can be triggered by external services.

Example: GitHub webhook triggers a deployment workflow

// List webhook workflows:
const webhooks = workflows.filter(w => w.trigger_type === 'webhook');

// Trigger via SDK (simulating webhook call):
await client.automation.workflows.trigger('wf_deploy_on_push', {
    source: 'github',
    event: 'push',
    repository: 'org/repo',
    branch: 'main',
    commit: 'abc123',
    author: 'developer@example.com'
});
`);
}

// ============================================
// Example: Workflow Orchestration
// ============================================

async function orchestrationExample() {
    console.log('\nWorkflow Orchestration Example\n');
    console.log('='.repeat(50));

    console.log(`
Complex business processes can be orchestrated using multiple workflows:

// 1. User signs up
const signupRun = await client.automation.workflows.trigger('wf_signup', {
    email: 'user@example.com',
    plan: 'enterprise'
});

// 2. Trigger provisioning workflow
const provisionRun = await client.automation.workflows.trigger('wf_provision', {
    userId: signupRun.id,
    resources: ['database', 'storage', 'compute']
});

// 3. Send welcome sequence
const welcomeRun = await client.automation.workflows.trigger('wf_welcome_sequence', {
    userId: signupRun.id,
    templateId: 'enterprise_welcome'
});

// 4. Notify sales team
await client.automation.workflows.trigger('wf_notify_sales', {
    userId: signupRun.id,
    plan: 'enterprise',
    priority: 'high'
});

console.log('Onboarding orchestration complete!');
`);
}

// ============================================
// Main
// ============================================

async function main() {
    try {
        await userOnboardingExample();
        await reportGenerationExample();
        await eventDrivenExample();
        await webhookExample();
        await orchestrationExample();
    } catch (error) {
        console.error('Error:', error);
        process.exit(1);
    }
}

main();
