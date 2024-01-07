# GCP Datastore Emulator UI

<!-- Badges -->

![License](https://img.shields.io/badge/license-MIT-blue)
![Version](https://img.shields.io/badge/version-0.0.1-blue)
![Go Version](https://img.shields.io/badge/go-1.21.5-blue)
![Next.js Version](https://img.shields.io/badge/next.js-14.0.4-blue)

A user-friendly, efficient, and modern UI for the GCP Datastore Emulator, designed for ease of use without compromising on resource efficiency.

<div style="display: flex; justify-content: center; align-items: center">
<img src="assets/gcp.png" height="40"/>
<img src="assets/datastore.png" height="60"/>
</div>

<!-- Table of Contents -->

## Table of Contents

- [Introduction](#introduction)
- [Key Features](#key-features)
- [Motivation](#motivation)
- [Getting Started](#getting-started)
- [Demo](#demo)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Support and Contact](#support-and-contact)

## Introduction

GCP Datastore Emulator UI is an open-source project aimed at providing a simple yet powerful user interface for the GCP Datastore Emulator. It's built with a focus on minimal resource usage, combining Golang and Next.js to serve a static site from the same server. This approach not only simplifies deployment but also ensures low overhead.

## Key Features

- **Lightweight Design**: Optimized for minimal resource consumption.
- **Basic Table View**: With type badges for easy comprehension.
- **Server-side Sorting**: Easily sort your data for better analysis.
- **Server-side Pagination**: Navigate through large datasets with ease.

## Motivation

This project is born out of the need for a practical and non-resource intensive UI for the GCP Datastore Emulator. While many solutions exist, they often compromise on either simplicity or efficiency. The goal with this implementation is to strike the perfect balance, with a keen eye on resource optimization. The ability to edit entities directly in the UI sets this project apart from existing solutions.

---

## Getting Started

Getting up and running with the GCP Datastore Emulator UI is straightforward. Follow these steps to build your own binaries using our `Makefile`.

### Building Your Own Binaries

1. **Clone the Repository**  
   Start by cloning this repository to your local machine.

   ```bash
    git clone https://github.com/Cyna298/gcp-datastore-ui.git
    cd gcp-datastore-ui
   ```

2. **Use the Makefile**

   We've provided a Makefile for easy building of the binaries. Make sure to update the `PROJECT_ID` variable in the Makefile to match your GCP project ID and `EMU_PORT` to match the port on which your emulator is running. You can also update the `PORT` variable to change the port on which the UI will be served.

   Navigate to the project's root directory and run:

   ```bash
   make setup
   ```

   This will compile the backend and frontend parts of the application and copy the frontend build to the backend's `out` directory.

3. **Run the Application**  
   After the build is complete, you can start the application using:

   ```bash
   make run-backend
   ```

   This will start the backend server and the frontend interface.

### Future Plans for Standalone Releases

We're working towards providing standalone releases of the GCP Datastore Emulator UI. These releases will offer a simplified, one-step setup process, making it easier than ever to get started. Stay tuned for updates!

---

## Demo

![Demo](assets/demo.gif)

## 🚀 Roadmap

We're on an exciting journey to enhance the GCP Datastore Emulator UI. Here's a glance at what we've achieved and what we're planning next:

- [x] **Basic Table View with Type Badges**  
       _Efficient display of data with clear type indication._

- [x] **Sorting**  
       _Organize and analyze your data with ease._

- [x] **Pagination**  
       _Navigate through large datasets more efficiently._

- [x] **Improved Build Process**  
       _Streamlining the build process for easier setup._

- [ ] **V0.1.0 Release**  
       _Our first major milestone on the horizon._

- [ ] **Edit Form / Detail View**  
       _Intuitive UI for directly editing entities._

- [ ] **Nested Entities**  
       _Handling complex data structures with ease._

- [ ] **Entity Navigation**  
       _Seamless navigation between related entities._

- [ ] **Simplified Table**  
       _Enhancing performance with a new table implementation._

---

## Contributing

We welcome contributions! Please open an issue or submit a pull request to propose changes.

## License

This project is licensed under the [MIT License](https://github.com/Cyna298/gcp-datastore-ui/blob/master/LICENSE).

## Support and Contact

For support, feature requests, or any queries, please [open an issue](https://github.com/Cyna298/gcp-datastore-ui/issues)

---
