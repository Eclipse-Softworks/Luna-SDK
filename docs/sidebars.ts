import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  // Manual sidebar configuration for optimal ordering
  tutorialSidebar: [
    'intro',
    'installation',
    'authentication',
    {
      type: 'category',
      label: 'Core Modules',
      items: [
        'payments',
        'messaging',
        'za-tools',
      ],
      collapsed: false,
    },
    'high-performance',
    {
      type: 'category',
      label: 'Reference',
      items: [
        'services',
        'errors',
        'versioning',
      ],
      collapsed: true,
    },
  ],
};

export default sidebars;
