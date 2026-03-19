from setuptools import setup, find_packages

setup(
    name="onelog-sdk",
    version="1.0.0",
    description="Official Python SDK for One-Log (ULAM) - Unified Log & Activity Monitor",
    author="One-Log Team",
    license="MIT",
    packages=find_packages(),
    install_requires=[
        "requests>=2.25.0",
    ],
    python_requires=">=3.7",
    keywords="logging monitoring apm observability ulam one-log",
    url="https://github.com/petrushandika/one-log",
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
)
