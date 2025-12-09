import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import { motion } from 'framer-motion';
import Link from '@docusaurus/Link';
import React from 'react';
import clsx from 'clsx';

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();

  return (
    <header className={clsx('hero', 'hero--primary')} style={{ background: 'transparent', textAlign: 'center' }}>
      <div className="container" style={{ maxWidth: '800px' }}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, ease: "easeOut" }}
        >
          <h1 className="hero__title">
            The Next Gen <br />
            <span style={{ color: '#fff', textShadow: '0 0 20px rgba(99,102,241,0.5)' }}>SDK Experience</span>
          </h1>
        </motion.div>

        <motion.p
          className="hero__subtitle"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
          Unified access to the Eclipse Softworks Platform. <br />
          Type-safe, robust, and designed for modern developers.
        </motion.p>

        <motion.div
          className="buttons"
          style={{ display: 'flex', gap: '20px', justifyContent: 'center', marginTop: '2rem' }}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.4 }}
        >
          <Link
            className="button button--primary button--lg"
            to="/docs/intro">
            Get Started
          </Link>
          <Link
            className="button button--secondary button--lg"
            to="https://github.com/Eclipse-Softworks/Luna-SDK">
            View on GitHub
          </Link>
        </motion.div>
      </div>
    </header>
  );
}

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - Modern SDKs`}
      description="Official SDKs for Eclipse Softworks Platform">

      {/* Background Glow Effect */}
      <div style={{
        position: 'absolute',
        top: '-10%',
        left: '50%',
        transform: 'translateX(-50%)',
        width: '600px',
        height: '600px',
        background: 'radial-gradient(circle, rgba(99,102,241,0.15) 0%, rgba(0,0,0,0) 70%)',
        zIndex: -1,
        pointerEvents: 'none'
      }} />

      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
