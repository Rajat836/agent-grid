import { SidebarLayout } from "@/components/sidebar-layout";

export default function Documentation() {
  return (
    <SidebarLayout>
      <div className="flex flex-1 flex-col gap-4 px-6 py-4 bg-gray-100 dark:bg-slate-950 h-full">
        {/* Header */}
        <div className="shrink-0 border-b border-gray-300 pb-4 rounded-lg p-4 bg-white dark:bg-slate-900">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            Documentation
          </h1>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1 font-medium">Explore comprehensive guides and resources</p>
        </div>

        {/* Content Area */}
        <div className="flex-1 rounded-lg border border-gray-300 p-8 bg-white dark:bg-slate-900 shadow-sm overflow-y-auto">
          <div className="space-y-6 max-w-4xl">
            <div className="p-6 rounded-lg bg-gray-50 dark:bg-slate-800 border border-gray-300 dark:border-slate-700">
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">Welcome to Documentation</h2>
              <p className="text-gray-700 dark:text-gray-300 leading-relaxed">
                This section contains comprehensive documentation to help you get started with the system analysis platform. Here you'll find guides, best practices, and resources to make the most of our tools.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="p-4 rounded-lg bg-gray-50 dark:bg-slate-800 border border-gray-300 dark:border-slate-700">
                <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">🚀 Getting Started</h3>
                <p className="text-sm text-gray-700 dark:text-gray-300">Begin your journey with step-by-step guides and tutorials to set up your workspace.</p>
              </div>
              <div className="p-4 rounded-lg bg-gray-50 dark:bg-slate-800 border border-gray-300 dark:border-slate-700">
                <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">📚 Guides & Tutorials</h3>
                <p className="text-sm text-gray-700 dark:text-gray-300">In-depth guides covering all major features and functionalities of the platform.</p>
              </div>
              <div className="p-4 rounded-lg bg-gray-50 dark:bg-slate-800 border border-gray-300 dark:border-slate-700">
                <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">⚙️ API Reference</h3>
                <p className="text-sm text-gray-700 dark:text-gray-300">Detailed API documentation for developers and integration specialists.</p>
              </div>
              <div className="p-4 rounded-lg bg-gray-50 dark:bg-slate-800 border border-gray-300 dark:border-slate-700">
                <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-2\">🔧 Troubleshooting</h3>
                <p className="text-sm text-gray-700 dark:text-gray-300">Common issues resolved with detailed solutions and workarounds.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </SidebarLayout>
  );
}
