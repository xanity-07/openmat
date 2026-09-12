import { ZLoginRequest, ZResponse } from '@openmat/zod';
import type { ZodOpenApiPathsObject } from 'zod-openapi';
import { getSecurityMetadata } from '../utils.js';

export const authPaths: ZodOpenApiPathsObject = {
    '/auth/login': {
        post: {
            tags: ['Auth'],
            summary: 'Login',
            description: 'authenticate with email and password, returns a JWT',
            requestBody: {
                content: {
                    'application/json': {
                        schema: ZLoginRequest,
                    },
                },
            },
            responses: {
                '200': {
                    description: 'login successful',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '400': {
                    description: 'invalid request',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '401': {
                    description: 'invalid email or password',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
            },
        },
    },
    '/auth/logout': {
        post: {
            tags: ['Auth'],
            summary: 'Logout',
            description: 'invalidate the current session',
            ...getSecurityMetadata({
                securityType: 'bearer',
            }),
            responses: {
                '204': {
                    description: 'logout successful',
                },
                '401': {
                    description: 'not authenticated or session already invalidated',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
            },
        },
    },
};
