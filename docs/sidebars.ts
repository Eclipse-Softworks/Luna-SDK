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
    {
      type: 'category',
      label: 'Getting Started',
      items: ['installation', 'authentication'],
    },
    {
      type: 'category',
      label: 'Core Concepts',
      items: [
        'core/users',
        'core/projects',
        'core/webhooks',
        'errors',
      ],
    },
    {
      type: 'category',
      label: 'High Performance',
      items: ['high-performance'],
    },
    {
      type: 'category',
      label: 'Cross-Platform',
      items: [
        'cross-platform/react-native',
        'cross-platform/electron',
        'cross-platform/pwa',
      ],
    },
    {
      type: 'category',
      label: 'Web Frameworks',
      items: [
        'web-frameworks/nextjs',
        'web-frameworks/vue-nuxt',
      ],
    },
    {
      type: 'category',
      label: 'Developer Tools',
      items: [
        'devtools/overview',
        'devtools/mock-client',
      ],
    },
    {
      type: 'category',
      label: 'SDK Reference',
      items: [
        {
          type: 'category',
          label: 'TypeScript',
          items: ['reference/typescript/client', 'reference/typescript/configuration'],
        },
        {
          type: 'category',
          label: 'Python',
          items: ['reference/python/client', 'reference/python/utilities'],
        },
        {
          type: 'category',
          label: 'Go',
          items: ['reference/go/client', 'reference/go/utilities'],
        },
      ],
    },
  ],
};

export default sidebars;
