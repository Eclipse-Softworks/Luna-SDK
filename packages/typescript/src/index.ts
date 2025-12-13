// Luna SDK - TypeScript Client
// Eclipse Softworks Platform API

export { LunaClient } from './client.js';
export type { ClientConfig, ClientOptions } from './client.js';

// Auth
export { ApiKeyAuth, TokenAuth } from './auth/index.js';
export type { AuthProvider, TokenPair, TokenRefreshCallback } from './auth/index.js';

// Resources
export { UsersResource } from './resources/users.js';
export { ProjectsResource } from './resources/projects.js';

// Errors
export {
    LunaError,
    AuthenticationError,
    AuthorizationError,
    ValidationError,
    RateLimitError,
    NetworkError,
    NotFoundError,
    ConflictError,
    ServerError,
} from './errors/index.js';
export { ErrorCode } from './errors/codes.js';

// Types
export type {
    User,
    UserCreate,
    UserUpdate,
    UserList,
    Project,
    ProjectCreate,
    ProjectUpdate,
    ProjectList,
    PaginationParams,
    ListResponse,
} from './types/index.js';

// HTTP
export type { RequestConfig, Response } from './http/index.js';

// Telemetry
export type { Logger, LogLevel, TelemetryConfig } from './telemetry/index.js';

// Payments (South African gateways)
export { Payments, PayFast, Ozow, Yoco, PayShap, SA_BANKS } from './resources/payments/index.js';
export type {
    PaymentsConfig,
    PayFastConfig,
    OzowConfig,
    YocoConfig,
    PayShapConfig,
    Payment,
    PayFastPayment,
    OzowPayment,
    YocoPayment,
    PayShapPayment,
    PaymentStatus,
    Currency,
    Amount,
} from './resources/payments/index.js';

// Messaging (SMS, WhatsApp, USSD)
export { Messaging, SMS, WhatsApp, USSD } from './resources/messaging/index.js';
export type {
    MessagingConfig,
    SMSConfig,
    WhatsAppConfig,
    USSDConfig,
    SMSMessage,
    WhatsAppMessage,
    USSDSession,
} from './resources/messaging/index.js';

// South African Business Tools
export { ZATools, CIPC, BBBEE, IDValidation, SAAddress } from './resources/za-tools/index.js';
export type {
    ZAToolsConfig,
    CIPCConfig,
    BBBEEConfig,
    Company,
    CompanyType,
    BBBEECertificate,
    BBBEELevel,
    SAIDInfo,
    Address,
    SAProvince,
} from './resources/za-tools/index.js';

