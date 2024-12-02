// src/app/(compliance)/privacy/page.tsx
import { ArrowUpRight } from "lucide-react";
import Link from "next/link";

const PrivacyPolicy = () => {
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
            <h1 className="text-4xl font-bold text-black">Privacy Policy</h1>
            <ArrowUpRight className="w-8 h-8 ml-2 text-black" />
          </div>

          <div className="prose prose-lg max-w-none">
            <p className="text-gray-600">Last updated: November 02, 2024</p>

            <div className="space-y-8">
              <section>
                <p className="text-gray-600 mt-4">
                  This Privacy Policy describes Our policies and procedures on
                  the collection, use and disclosure of Your information when
                  You use the Service and tells You about Your privacy rights
                  and how the law protects You.
                </p>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">
                  Types of Data Collected
                </h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <h3 className="text-xl font-semibold mb-2">Personal Data</h3>
                  <p className="text-gray-600 mb-4">
                    While using Our Service, We may ask You to provide Us with
                    certain personally identifiable information that can be used
                    to contact or identify You. Personally identifiable
                    information may include, but is not limited to:
                  </p>
                  <ul className="list-disc list-inside text-gray-600 space-y-2">
                    <li>Email address</li>
                    <li>Usage Data</li>
                  </ul>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">
                  Use of Your Personal Data
                </h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600 mb-4">
                    The Company may use Personal Data for the following
                    purposes:
                  </p>
                  <ul className="list-disc list-inside text-gray-600 space-y-2">
                    <li>To provide and maintain our Service</li>
                    <li>To manage Your Account</li>
                    <li>For the performance of a contract</li>
                    <li>To contact You</li>
                    <li>To manage Your requests</li>
                  </ul>
                </div>
              </section>

              <section>
                <h2 className="text-2xl font-bold mt-8 mb-4">Contact Us</h2>
                <div className="pl-4 border-l-2 border-gray-200">
                  <p className="text-gray-600">
                    If you have any questions about this Privacy Policy, You can
                    contact us:
                  </p>
                  <p className="text-gray-600 mt-2">
                    By email:{" "}
                    <a
                      href="mailto:hey@byrdhq.com"
                      className="text-blue-600 hover:underline"
                    >
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

export default PrivacyPolicy;
