import Link from 'next/link';
import { ArrowUpRight } from 'lucide-react';

const TermsOfService = () => {
  return (
    <div className="min-h-screen bg-white text-black flex flex-col">
      <nav>
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <Link href="/" className="text-xl font-bold">
            byrd
          </Link>
        </div>
      </nav>

      <main className="flex-grow px-4 py-12">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center mb-8">
            <h1 className="text-4xl font-bold text-black">Terms of Service</h1>
            <ArrowUpRight className="w-8 h-8 ml-2 text-black" />
          </div>

          <div className="prose prose-lg max-w-none">
            <p className="text-gray-600">Last updated: November 03, 2024</p>

            <div className="space-y-8">
              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">1. Acceptance of Terms</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    By accessing or using byrd's services, you agree to be bound by these Terms of Service. If you disagree
                    with any part of the terms, you may not access the service.
                  </p>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">2. Description of Service</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    byrd provides a software-as-a-service platform for competitive intelligence gathering and analysis.
                    The specific features and functionality may be modified, updated, or removed at any time at our discretion.
                  </p>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">3. User Accounts</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600 mb-4">
                    To access the service, you must create an account. You agree to:
                  </p>
                  <ul className="list-disc list-inside text-gray-600 space-y-2">
                    <li>Provide accurate account information</li>
                    <li>Maintain the security of your account</li>
                    <li>Promptly update any changes to your account information</li>
                    <li>Accept responsibility for all activities under your account</li>
                  </ul>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">4. Service Usage</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600 mb-4">Users agree not to:</p>
                  <ul className="list-disc list-inside text-gray-600 space-y-2">
                    <li>Use the service for any illegal purposes</li>
                    <li>Share account credentials with unauthorized users</li>
                    <li>Attempt to circumvent any usage limitations</li>
                    <li>Reverse engineer the service</li>
                    <li>Use the service in a way that could damage or overburden our systems</li>
                  </ul>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">5. Data and Privacy</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    Our collection and use of personal information is governed by our Privacy Policy. By using the service,
                    you consent to our data practices as described in the Privacy Policy.
                  </p>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">6. Intellectual Property</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    The service, including all related features, functionality, and content, is owned by byrd and protected
                    by intellectual property laws. Users retain ownership of their data but grant us necessary licenses to
                    provide the service.
                  </p>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">7. Disclaimers</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600 mb-4">
                    The service is provided "as is" without warranties of any kind. We do not guarantee:
                  </p>
                  <ul className="list-disc list-inside text-gray-600 space-y-2">
                    <li>Continuous, uninterrupted access to the service</li>
                    <li>The accuracy or completeness of any data or analysis</li>
                    <li>The preservation or backup of user data</li>
                    <li>That the service will meet your specific requirements</li>
                  </ul>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">8. Limitation of Liability</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    To the maximum extent permitted by law, byrd shall not be liable for any indirect, incidental, special,
                    consequential, or punitive damages, or any loss of profits or revenues.
                  </p>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">9. Changes to Terms</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    We reserve the right to modify these terms at any time. We will notify users of any material changes via
                    email or through the service. Continued use of the service after such modifications constitutes acceptance
                    of the updated terms.
                  </p>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">10. Contact</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    For questions about these Terms, please contact us at:{' '}
                    <a href="mailto:hey@byrdhq.com" className="text-blue-600 hover:underline">
                      hey@byrdhq.com
                    </a>
                  </p>
                </div>
              </section>
            </div>
          </div>
        </div>
      </main>

      <footer className="py-8 px-6">
        <div className="max-w-7xl mx-auto">
          <div className="flex justify-between items-center text-sm text-gray-600">
            <p>© 2024 byrd. All rights reserved.</p>
            <div className="flex space-x-4">
              <Link href="/terms" className="hover:text-black">
                Terms
              </Link>
              <Link href="/privacy" className="hover:text-black">
                Privacy
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default TermsOfService;
