"""
Setup script for BTNG SDK
"""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="btng_sdk",
    version="0.1.0",
    author="BTNG Sovereign Platform",
    author_email="contact@btng.africa",
    description="Unified Python SDK for BTNG blockchain operations",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/btng-sovereign/btng-sdk",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    python_requires=">=3.8",
    install_requires=[
        "ecdsa>=0.18.0",
        "aiohttp>=3.8.0",
        "requests>=2.28.0",
    ],
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "black>=22.0.0",
            "flake8>=4.0.0",
        ],
    },
)