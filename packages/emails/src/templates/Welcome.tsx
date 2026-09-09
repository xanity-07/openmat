import {
    Body,
    Button,
    Container,
    Head,
    Heading,
    Hr,
    Link,
    Preview,
    Section,
    Tailwind,
    Text
} from "@react-email/components";

interface WelcomeEmailProps {
    userFirstName: string;
}

export const WelcomeEmail = ({ userFirstName = "{{.UserFirstName}}" }: WelcomeEmailProps) => {
    return (
        <html>
        <Tailwind>
            <Head>
                <style>{`
                        @media (prefers-color-scheme: light) {
                            .om-body { background-color: #f1f5f9 !important; }
                            .om-container { background-color: #ffffff !important; }
                            .om-heading { color: #1e293b !important; }
                            .om-text { color: #334155 !important; }
                            .om-hr { border-color: #e2e8f0 !important; }
                            .om-link { color: #dc2626 !important; }
                            .om-footer { color: #64748b !important; }
                        }
                    `}</style>
            </Head>
            <Preview>Welcome To OpenMat</Preview>
            <Body className="om-body bg-black font-sans">
                <Container className="om-container bg-neutral-900 p-8 rounded-lg shadow-sm my-10 mx-auto max-w-[600px]">
                    <Heading className='om-heading text-2xl font-bold text-white mt-4'>
                        Welcome To OpenMat!!
                    </Heading>

                    <Section>
                        <Text className='om-text text-neutral-300 text-base'>
                            Hi {userFirstName}
                        </Text>
                        <Text className='om-text text-neutral-300 text-base'>
                            Thank you for joining !!
                        </Text>
                    </Section>

                    <Section className='my-8 text-center'>
                        <Button
                            className='bg-red-600 hover:bg-red-900 text-white font-medium rounded-md px-6 py-3'
                            href={`/dashboard`}
                        >
                            Get Started
                        </Button>
                    </Section>

                    <Hr className='om-hr border-neutral-700 my-6' />

                    <Section>
                        <Text className='om-text text-neutral-300 text-base'>
                            If you have any questions, feel free to {" "}
                            <Link href={`/support`} className='om-link text-red-400 underline'>
                                contact our support team
                            </Link>
                            .
                        </Text>
                    </Section>

                    <Section className='mt-8 text-center'>
                        <Text className='om-footer text-neutral-500 text-xs'>
                            © {new Date().getFullYear()} OpenMat. All rights reserved.
                        </Text>
                        <Text className='om-footer text-neutral-500 text-xs'>
                            123 Project Street, Suite 100, San Francisco, CA 94103
                        </Text>
                    </Section>
                </Container>
            </Body>
        </Tailwind>
        </html>
    )
}

WelcomeEmail.PreviewProps = {
    userFirstName: "John"
}

export default WelcomeEmail