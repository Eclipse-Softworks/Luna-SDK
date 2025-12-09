import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Luna SDK',
  tagline: 'The Official Eclipse Softworks Platform SDK',
  favicon: 'img/favicon.ico',

  // Production URL
  url: 'https://docs.eclipse-softworks.com',
  // Sub-path for the site
  baseUrl: '/luna-sdk/',

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
    image: 'img/docusaurus-social-card.jpg',
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
