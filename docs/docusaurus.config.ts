import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Luna SDK',
  tagline: 'The Official Eclipse Softworks Platform SDK',
  favicon: 'img/favicon.ico',

  // Production URL
  url: 'https://docs-lunasdk.eclipse-softworks.com',
  // Sub-path for the site
  baseUrl: '/',

  // GitHub config
  organizationName: 'Eclipse-Softworks',
  projectName: 'Luna-SDK',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl:
            'https://github.com/Eclipse-Softworks/Luna-SDK/tree/master/docs/',
        },
        blog: false, // Disable blog for now to keep it clean
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/social-card.png',
    metadata: [
      { name: 'description', content: 'Official TypeScript, Python, and Go SDKs for the Eclipse Softworks Platform. Type-safe, production-ready client libraries.' },
      { name: 'keywords', content: 'Luna SDK, Eclipse Softworks, TypeScript SDK, Python SDK, Go SDK, API Client' },
      { property: 'og:title', content: 'Luna SDK - Eclipse Softworks' },
      { property: 'og:description', content: 'Official multi-language SDKs for the Eclipse Softworks Platform API.' },
      { property: 'og:type', content: 'website' },
      { property: 'og:url', content: 'https://docs-lunasdk.eclipse-softworks.com' },
      { property: 'og:image', content: 'https://docs-lunasdk.eclipse-softworks.com/img/social-card.png' },
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:title', content: 'Luna SDK - Eclipse Softworks' },
      { name: 'twitter:description', content: 'Official SDKs for TypeScript, Python, and Go.' },
      { name: 'twitter:image', content: 'https://docs-lunasdk.eclipse-softworks.com/img/social-card.png' },
    ],
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Luna SDK',
      logo: {
        alt: 'Luna SDK Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Documentation',
        },
        {
          href: 'https://github.com/Eclipse-Softworks/Luna-SDK',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/intro',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/Eclipse-Softworks/Luna-SDK',
            },
            {
              label: 'Eclipse Softworks',
              href: 'https://www.eclipse-softworks.com',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Eclipse Softworks. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'json', 'typescript', 'python', 'go'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
