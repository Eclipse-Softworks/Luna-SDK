import React from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import { motion } from 'framer-motion';
import { SiTypescript, SiPython, SiGo } from 'react-icons/si';
import { FiShield, FiZap, FiGlobe } from 'react-icons/fi';

type FeatureItem = {
  title: string;
  icon: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Multi-Language SDKs',
    icon: SiTypescript, // Representing the suite
    description: (
      <>
        First-class support for <strong>TypeScript</strong>, <strong>Python</strong>, and <strong>Go</strong>.
        Write idiomatic code in your language of choice with fully typed clients.
      </>
    ),
  },
  {
    title: 'Production Ready',
    icon: FiShield,
    description: (
      <>
        Built for reliability with automatic <strong>retries</strong>, <strong>timeout management</strong>,
        and robust error handling out of the box. Enterprise-grade stability.
      </>
    ),
  },
  {
    title: 'Secure by Default',
    icon: FiZap,
    description: (
      <>
        Integrated <strong>OAuth</strong> and <strong>API Key</strong> authentication strategies.
        Secure storage helpers for tokens and automatic rotation handling.
      </>
    ),
  },
];

function Feature({ title, icon: Icon, description }: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="feature-card margin-bottom--lg">
        <div className="feature-icon">
          <Icon />
        </div>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features} style={{ padding: '4rem 0' }}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
