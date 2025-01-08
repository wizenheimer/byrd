"use client";
import { ChevronDown, Search } from 'lucide-react';
import { useState } from 'react';
import { UserButton } from '@clerk/nextjs';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger
} from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';

export default function Dashboard() {
  // Add state for tracking selected item and group
  const [activeTab, setActiveTab] = useState('history');
  const [selectedItem, setSelectedItem] = useState<{ name: string; url: string } | null>(null);
  const [selectedGroup, setSelectedGroup] = useState<string | null>(null);
  const [groupVisibility, setGroupVisibility] = useState<{ [key: string]: boolean }>({});

  const sidebarItems = [
    {
      title: 'Common Room',
      items: Array(5).fill({
        name: 'Common Room Signal',
        url: 'https://www.commonroom.io/',
      })
    },
    {
      title: 'Reo Dev',
      items: Array(5).fill({
        name: 'Reo Dev',
        url: 'https://www.reo.dev/',
      })
    },
    {
      title: 'Scarf Analytics',
      items: Array(5).fill({
        name: 'Scarf Analytics',
        url: 'https://about.scarf.sh/',
      })
    }
  ];

  // Define changes items
  const changesItems = [
    {
      title: 'Branding',
      details: []
    },
    {
      title: 'Pricing',
      details: [
        'Increased premium plan pricing by 10%',
        'Added a new enterprise tier',
        'Removed free tier'
      ]
    },
    {
      title: 'Positioning',
      details: [
        'Updated homepage headline',
        'Added new customer testimonials',
        'Updated pricing page copy',
        'Added new case studies'
      ]
    },
    {
      title: 'Features',
      details: []
    },
    {
      title: 'Channels',
      details: [
        "Added new channel: TikTok",
        "Removed channel: Pinterest"
      ]
    },
    {
      title: 'Team',
      details: []
    }
  ];

  const historyItems = [
    { time: 'in 10 hours', status: 'Scheduled' },
    { time: '2 days ago', status: 'Review' },
    { time: '3 days ago', status: 'Review' },
    { time: '4 days ago', status: 'Review' }
  ];

  // Toggle group visibility
  const toggleGroup = (groupTitle: string) => {
    setGroupVisibility(prev => ({
      ...prev,
      [groupTitle]: !prev[groupTitle]
    }));
    setSelectedGroup(groupTitle);
  };

  // Handle item selection
  const handleItemSelect = (item: { name: string; url: string }, group: { title: string }) => {
    setSelectedItem(item);
    setSelectedGroup(group.title);
  };

  return (
    <div className="h-screen w-screen flex bg-white">
      {/* Left Sidebar */}
      <div className="w-[420px] flex flex-col border-r border-gray-200">
        <div className="h-[56px] flex items-center px-4">
          <span className="text-xl font-bold">byrd</span>
        </div>

        <div className="p-4 flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-400" />
            <Input
              placeholder="Search Pages"
              className="pl-8 h-9 text-sm bg-transparent"
            />
          </div>
          <Button
            variant="outline"
            className="h-9 w-9 p-0 flex items-center justify-center bg-[#171717] text-white rounded-lg"
          >
            <span className="text-lg">+</span>
          </Button>
        </div>

        <div className="px-4 py-2">
          <span className="text-xs text-gray-500">Active Pages (15)</span>
        </div>

        <div className="flex-1 overflow-y-auto">
          {sidebarItems.map((group, groupIndex) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
            <div key={groupIndex}>
              <button
                type="button"
                className="w-full px-4 py-2 flex items-center gap-1 text-sm font-medium"
                onClick={() => toggleGroup(group.title)}
              >
                <svg
                  className={`w-3 h-3 text-gray-500 transform transition-transform ${groupVisibility[group.title] ? 'rotate-180' : ''
                    }`}
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                >
                  <title>{ }</title>
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                </svg>
                {group.title}
              </button>
              {groupVisibility[group.title] && group.items.map((item, itemIndex) => (
                // biome-ignore lint/a11y/useKeyWithClickEvents: <explanation>
                <div
                  // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
                  key={itemIndex}
                  className={`flex items-center px-4 py-2 hover:bg-gray-50 group cursor-pointer ${selectedItem === item && selectedGroup === group.title ? 'bg-gray-50' : ''
                    }`}
                  onClick={() => handleItemSelect(item, group)}
                >
                  <img
                    src={`https://www.google.com/s2/favicons?domain=${new URL(item.url).hostname}&sz=32`}
                    alt=""
                    className="w-4 h-4 mr-2"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="text-sm text-gray-900">{item.name}</div>
                    <div className="text-xs text-gray-500 truncate">{item.url}</div>
                  </div>
                  <button className="opacity-0 group-hover:opacity-100 ml-2" type="button">
                    <svg className="w-4 h-4 text-gray-400" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                      <title>{ }</title>
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          ))}
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col">
        <div className="h-[56px] flex items-center justify-between px-4 border-b border-gray-200">
          <div className="flex items-center gap-2">
            {selectedItem && (
              <>
                <img
                  src={`https://www.google.com/s2/favicons?domain=${new URL(selectedItem.url).hostname}&sz=32`}
                  alt=""
                  className="w-4 h-4"
                />
                <span className="text-sm font-medium">{selectedItem.name}</span>
              </>
            )}
          </div>
          <div className="flex items-center gap-2">
            <button type="button">
              <svg className="w-4 h-4 text-gray-400" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <title>{ }</title>
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
              </svg>
            </button>
            <UserButton />
          </div>
        </div>

        <div className="flex-1 p-6">
          {selectedItem ? (
            <>
              <div className="w-full h-[400px] border border-gray-100 rounded-lg" />
              <div className="mt-6">
                <Tabs value={activeTab} onValueChange={setActiveTab}>
                  <div className="flex justify-center">
                    <TabsList className="w-[200px]">
                      <TabsTrigger value="history">History</TabsTrigger>
                      <TabsTrigger value="changes">Changes</TabsTrigger>
                    </TabsList>
                  </div>

                  <TabsContent value="history" className="mt-4">
                    {historyItems.map((item, index) => (
                      <div
                        // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
                        key={index}
                        className="flex items-center justify-between p-4 border border-gray-100 rounded-lg mb-2"
                      >
                        <div className="flex items-center gap-2">
                          <svg className="w-4 h-4 text-gray-400" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                            <title>{ }</title>
                            <circle cx="12" cy="12" r="10" strokeWidth="2" />
                            <path d="M12 6v6l4 2" strokeWidth="2" />
                          </svg>
                          <span className="text-sm text-gray-500">{item.time}</span>
                        </div>
                        <span className="text-sm text-gray-500">{item.status}</span>
                      </div>
                    ))}
                  </TabsContent>

                  <TabsContent value="changes" className="mt-4 space-y-3">
                    {changesItems.map((item, index) => (
                      // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
                      <Collapsible key={index} className="border rounded-lg">
                        <CollapsibleTrigger className="flex items-center justify-between w-full p-4 text-left">
                          <div className="flex items-center gap-3">
                            <div className="font-medium">{item.title}</div>
                            {item.details.length > 0 && (
                              <div className="px-2 py-0.5 bg-gray-100 rounded-full text-sm text-gray-600">
                                {item.details.length} Changes
                              </div>
                            )}
                          </div>
                          <div className="flex items-center gap-4">
                            {item.details.length > 0 && (
                              <ChevronDown className="h-4 w-4 text-gray-500 transition-transform duration-200 data-[state=open]:rotate-180" />
                            )}
                          </div>
                        </CollapsibleTrigger>
                        {item.details.length > 0 && (
                          <CollapsibleContent className="px-4 pb-4">
                            <ul className="text-sm text-gray-600 space-y-2 list-disc ml-4">
                              {item.details.map((detail, detailIndex) => (
                                // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
                                <li key={detailIndex}>{detail}</li>
                              ))}
                            </ul>
                          </CollapsibleContent>
                        )}
                      </Collapsible>
                    ))}
                  </TabsContent>
                </Tabs>
              </div>
            </>
          ) : (
            <div className="flex items-center justify-center h-full text-gray-500">
              Select an item from the sidebar to view details
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
