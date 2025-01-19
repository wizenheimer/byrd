import { ArrowUpRight } from "lucide-react";
import Link from "next/link";

const EULA = () => {
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
						<h1 className="text-4xl font-bold text-black">
							End User License Agreement
						</h1>
						<ArrowUpRight className="w-8 h-8 ml-2 text-black" />
					</div>

					<div className="prose prose-lg max-w-none space-y-8">
						<p className="text-gray-600">Last Updated: November 03, 2024</p>

						<section>
							<p className="text-gray-600">
								This End User License Agreement (&quot;Agreement&quot;) is a
								legal agreement between you (&quot;User&quot;, &quot;you&quot;,
								or &quot;your&quot;) and ByrdLabs, Inc. (&quot;Company&quot;,
								&quot;we&quot;, &quot;us&quot;, or &quot;our&quot;) for the use
								of our competitive intelligence software-as-a-service platform
								(&quot;Service&quot;).
							</p>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">1. License Grant</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									Subject to the terms of this Agreement and payment of
									applicable fees, we grant you a limited, non-exclusive,
									non-transferable, non-sublicensable license to:
								</p>
								<ul className="list-disc list-inside text-gray-600 mt-4 space-y-2">
									<li>
										Access and use the Service for your internal business
										purposes
									</li>
									<li>Generate and download reports using the Service</li>
									<li>Store and process data through the Service</li>
								</ul>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">2. Restrictions</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">You agree not to:</p>
								<ul className="list-disc list-inside text-gray-600 mt-4 space-y-2">
									<li>Resell, sublicense, or distribute the Service</li>
									<li>Modify, reverse engineer, or create derivative works</li>
									<li>Remove any proprietary notices or labels</li>
									<li>Use the Service to build a competitive product</li>
									<li>Exceed usage limits or circumvent access restrictions</li>
								</ul>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">3. Ownership</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									The Service, including all intellectual property rights,
									belongs to us. You retain ownership of any data you input into
									the Service. You grant us a license to use your data to
									provide and improve the Service.
								</p>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">
								4. Confidentiality
							</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									You agree to maintain the confidentiality of any non-public
									features, functionality, or information about the Service.
									This includes beta features and pre-release information.
								</p>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">
								5. Term and Termination
							</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									This Agreement remains in effect until terminated. We may
									suspend or terminate your access:
								</p>
								<ul className="list-disc list-inside text-gray-600 mt-4 space-y-2">
									<li>For violation of this Agreement</li>
									<li>For non-payment of fees</li>
									<li>If required by law</li>
									<li>Upon discontinuation of the Service</li>
								</ul>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">
								6. Warranty Disclaimer
							</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									THE SERVICE IS PROVIDED &quot;AS IS&quot; WITHOUT WARRANTY OF
									ANY KIND. WE DISCLAIM ALL WARRANTIES, WHETHER EXPRESS,
									IMPLIED, OR STATUTORY, INCLUDING MERCHANTABILITY, FITNESS FOR
									A PARTICULAR PURPOSE, AND NON-INFRINGEMENT.
								</p>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">
								7. Limitation of Liability
							</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									TO THE MAXIMUM EXTENT PERMITTED BY LAW, WE SHALL NOT BE LIABLE
									FOR ANY INDIRECT, INCIDENTAL, SPECIAL, CONSEQUENTIAL, OR
									PUNITIVE DAMAGES, OR ANY LOSS OF PROFITS OR REVENUES.
								</p>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">
								8. Indemnification
							</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									You agree to indemnify and hold us harmless from any claims
									arising from your use of the Service or violation of this
									Agreement.
								</p>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">9. Modifications</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									We reserve the right to modify this Agreement at any time. We
									will notify you of material changes. Continued use of the
									Service constitutes acceptance of modifications.
								</p>
							</div>
						</section>

						<section>
							<h2 className="text-2xl font-bold mt-8 mb-4">
								10. Governing Law
							</h2>
							<div className="pl-4 border-l-2 border-gray-200">
								<p className="text-gray-600">
									This Agreement is governed by the laws of soverign republic
									state of India, without regard to its conflict of law
									principles.
								</p>
							</div>
						</section>

						<section className="mt-8">
							<p className="text-gray-600">
								By using the Service, you acknowledge that you have read,
								understood, and agree to be bound by this Agreement.
							</p>
						</section>
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

export default EULA;
