/**
 * Luna SDK - Messaging Module Types
 * SMS, WhatsApp, USSD, Email integrations for Africa
 */

// ============================================================================
// Common Types
// ============================================================================

export type MessageStatus =
    | 'pending'
    | 'sent'
    | 'delivered'
    | 'read'
    | 'failed'
    | 'expired';

export type MessageDirection = 'outbound' | 'inbound';

export interface MessageMetadata {
    [key: string]: string | number | boolean;
}

// ============================================================================
// SMS Types
// ============================================================================

export type SMSProvider = 'clickatell' | 'africastalking' | 'twilio';

export interface SMSConfig {
    provider: SMSProvider;
    apiKey: string;
    apiSecret?: string;
    /** Sender ID or short code */
    senderId?: string;
    /** Africa's Talking username */
    username?: string;
    /** Twilio Account SID */
    accountSid?: string;
    sandbox?: boolean;
}

export interface SMSMessage {
    id: string;
    to: string;
    from?: string;
    body: string;
    status: MessageStatus;
    direction: MessageDirection;
    provider: SMSProvider;
    /** Cost in ZAR cents */
    cost?: number;
    /** Number of message parts */
    parts?: number;
    metadata?: MessageMetadata;
    created_at: string;
    updated_at: string;
}

export interface SMSSendRequest {
    to: string | string[];
    body: string;
    from?: string;
    /** Schedule for later delivery (ISO 8601) */
    scheduledAt?: string;
    /** Webhook URL for delivery reports */
    callbackUrl?: string;
    metadata?: MessageMetadata;
}

export interface SMSBulkResult {
    successful: SMSMessage[];
    failed: Array<{ to: string; error: string }>;
}

// ============================================================================
// WhatsApp Types
// ============================================================================

export type WhatsAppProvider = 'cloud_api' | 'clickatell' | 'wati' | 'infobip';

export interface WhatsAppConfig {
    provider: WhatsAppProvider;
    /** Business Account ID or API Key */
    apiKey: string;
    /** Phone number ID (Cloud API) */
    phoneNumberId?: string;
    /** Access token for Cloud API */
    accessToken?: string;
    /** Webhook verify token */
    webhookToken?: string;
    sandbox?: boolean;
}

export type WhatsAppMessageType =
    | 'text'
    | 'template'
    | 'image'
    | 'document'
    | 'audio'
    | 'video'
    | 'location'
    | 'contacts'
    | 'interactive';

export interface WhatsAppMessage {
    id: string;
    to: string;
    from?: string;
    type: WhatsAppMessageType;
    status: MessageStatus;
    direction: MessageDirection;
    provider: WhatsAppProvider;
    /** Text content for text messages */
    text?: string;
    /** Template name for template messages */
    templateName?: string;
    /** Template parameters */
    templateParams?: Record<string, string>;
    /** Media URL for media messages */
    mediaUrl?: string;
    /** Conversation ID within 24-hour window */
    conversationId?: string;
    metadata?: MessageMetadata;
    created_at: string;
    updated_at: string;
}

export interface WhatsAppTextRequest {
    to: string;
    text: string;
    previewUrl?: boolean;
}

export interface WhatsAppTemplateRequest {
    to: string;
    templateName: string;
    templateParams?: Record<string, string>;
    language?: string;
}

export interface WhatsAppMediaRequest {
    to: string;
    type: 'image' | 'document' | 'audio' | 'video';
    mediaUrl: string;
    caption?: string;
    filename?: string;
}

export interface WhatsAppWebhookPayload {
    object: string;
    entry: Array<{
        id: string;
        changes: Array<{
            value: {
                messaging_product: string;
                metadata: { phone_number_id: string };
                messages?: Array<{
                    id: string;
                    from: string;
                    type: string;
                    text?: { body: string };
                    timestamp: string;
                }>;
                statuses?: Array<{
                    id: string;
                    status: string;
                    timestamp: string;
                }>;
            };
        }>;
    }>;
}

// ============================================================================
// USSD Types
// ============================================================================

export interface USSDConfig {
    provider: 'clickatell' | 'africastalking';
    apiKey: string;
    /** USSD service code (e.g., *120*123#) */
    serviceCode: string;
    username?: string;
    sandbox?: boolean;
}

export type USSDSessionState = 'active' | 'ended' | 'timeout';

export interface USSDSession {
    id: string;
    sessionId: string;
    phoneNumber: string;
    serviceCode: string;
    text: string;
    state: USSDSessionState;
    /** Network: MTN, Vodacom, Cell C, Telkom */
    network?: string;
    created_at: string;
    updated_at: string;
}

export interface USSDResponse {
    /** Response text to display */
    text: string;
    /** Whether to end the session */
    end: boolean;
}

export interface USSDHandler {
    (session: USSDSession): Promise<USSDResponse> | USSDResponse;
}

// ============================================================================
// Email Types (optional extension)
// ============================================================================

export interface EmailConfig {
    provider: 'sendgrid' | 'mailgun' | 'ses';
    apiKey: string;
    domain?: string;
    fromEmail: string;
    fromName?: string;
}

export interface EmailMessage {
    id: string;
    to: string | string[];
    from: string;
    subject: string;
    body: string;
    bodyHtml?: string;
    status: MessageStatus;
    created_at: string;
}

export interface EmailSendRequest {
    to: string | string[];
    subject: string;
    body: string;
    bodyHtml?: string;
    from?: string;
    replyTo?: string;
    cc?: string[];
    bcc?: string[];
    attachments?: Array<{
        filename: string;
        content: string | Buffer;
        contentType?: string;
    }>;
}

// ============================================================================
// Delivery Report Types
// ============================================================================

export interface DeliveryReport {
    messageId: string;
    status: MessageStatus;
    deliveredAt?: string;
    failureReason?: string;
    provider: string;
}

export interface DeliveryReportWebhook {
    (report: DeliveryReport): Promise<void> | void;
}
