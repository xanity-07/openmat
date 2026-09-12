import { userPaths } from '@/contracts/user.js';
import { createDocument, type ZodOpenApiObject } from 'zod-openapi';
import { authPaths } from './contracts/auth.js';
import { healthPaths } from './contracts/health.js';

const openApiConfig: ZodOpenApiObject = {
    openapi: '3.1.0',
    info: {
        version: '1.0.0',
        title: 'OpenMat REST API - Documentation',
        description: 'OpenMat REST API - Documentation',
    },
    servers: [
        {
            url: 'http://localhost:8080/',
            description: 'Local Server',
        },
        {
            url: 'http://localhost:8080/api/v1',
            description: 'Local Server',
        },
    ],
    paths: {
        ...healthPaths,
        ...authPaths,
        ...userPaths,
    },
    components: {
        schemas: {},
        securitySchemes: {
            bearerAuth: {
                type: 'http',
                scheme: 'bearer',
                bearerFormat: 'JWT',
            },
            'x-service-token': {
                type: 'apiKey',
                name: 'x-service-token',
                in: 'header',
            },
        },
    },
};

export const OpenAPI: ReturnType<typeof createDocument> = createDocument(openApiConfig);
