<img src="assets/gcp.png" height="50" align="right"/>
<!-- Badges -->

![License](https://img.shields.io/badge/license-MIT-blue)
![Version](https://img.shields.io/badge/version-0.0.2-blue)
![Go Version](https://img.shields.io/badge/go-1.21.5-blue)
![Templ Version](https://img.shields.io/badge/templ-latest-blue)
![Htmx Version](https://img.shields.io/badge/htmx-latest-blue)

A user-friendly, efficient, and modern UI for the GCP Datastore Emulator, designed for ease of use without compromising on resource efficiency. This project leans on server-side rendering to simplify the architecture and focus on minimalistic and efficient design.

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

GCP Datastore Emulator UI is an open-source project aimed at providing a simple yet powerful user interface for the GCP Datastore Emulator. It's built with a focus on minimal resource usage, combining `Golang`, `Templ`, and `HTMX` to serve a server-side rendered application. This approach ensures low overhead and adheres to the principles of minimalism and resource efficiency.

## Key Features

- **Lightweight and Server-side Rendered Design**: Focused on minimal resource consumption without reliance on client-side frameworks.
- **Server-side Sorting and Pagination**: Navigate and manage large datasets efficiently.
- **MVVM Architecture Inspiration**: Maintaining all state on the backend to simplify the client-side as a pure view representation.

## Motivation

This project is designed to address the need for an efficient, yet simple UI for the GCP Datastore Emulator that does not sacrifice performance for simplicity. By moving away from traditional web app frameworks and embracing a more desktop-app-like architecture, it offers a unique approach to managing and interacting with datastore entities.

## Getting Started

Follow these steps to set up the GCP Datastore Emulator UI:

### Building Your Own Binaries

1. **Clone the Repository**  
   Start by cloning this repository to your local machine.

   ```bash
    git clone https://github.com/Cyna298/gcp-datastore-ui.git
    cd gcp-datastore-ui
   ```

2. **Install Dependencies**

   A `Makefile` will be provided in the future. For now download `templ`, `tailwind-cli` and `air`

3. **Run the Application**  
   Initiate the application using:

   ```bash
   make run-backend
   ```

   This launches the backend server, rendering the frontend through server-side templating.

### Future Plans

- **TUI Interface**: Exploring a terminal user interface to completely move away from the web aspect.

## Demo

![Demo](assets/demo.gif)

## Roadmap

Future developments and current features include:

- [ ] **Basic Table View with Type Badges**  
       _Efficient display of data with clear type indication._

- [x] **Sorting**  
       _Organize and analyze your data with ease._

- [x] **Pagination**  
       _Navigate through large datasets more efficiently._

- [ ] **Improved Build Process**  
       _Streamlining the build process for easier setup._

- [x] **V0.1.0 Release**  
       _Our first major milestone on the horizon._

- [ ] **Edit Form / Detail View**  
       _Intuitive UI for directly editing entities._

- [x] **Nested Entities**  
       _Handling complex data structures with ease._

- [ ] **Entity Navigation**  
       _Seamless navigation between related entities._

- [x] **Simplified Table**  
       _Enhancing performance with a new table implementation._

- [ ] **TUI Based Interface**
      Experimenting with BubbleTea for a text-based user interface.
- [x] **MVVM Architecture**:
      All state managed on the backend, with the frontend as a reflection of this state.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request to suggest changes.

## License

This project is licensed under the [MIT License](https://github.com/Cyna298/gcp-datastore-ui/blob/master/LICENSE).

## Support and Contact

For support, feature requests, or queries, please [open an issue](https://github.com/Cyna298/gcp-datastore-ui/issues).
