import os
import re
import yaml

def extract_from_adr(adr_content, module_name):
    module_data = {
        "title": None,
        "description": None,
        "criticality": "MEDIUM", # Default
        "implementation_phase": None,
        "key_entities": [],
        "dependencies_on": [],
        "depended_on_by": [],
        "use_cases": [],
        "api_endpoints": []
    }

    # Extract Title (from H1)
    title_match = re.search(r"# ADR-\d+ – (.+)", adr_content)
    if title_match:
        module_data["title"] = title_match.group(1).strip()
    else: # Fallback to ADR file name as title
        module_data["title"] = module_name.replace("-", " ").title()

    # Extract description (often in Context section)
    desc_match = re.search(r"## 1\. Contexto\n\n(.+?)\n\n---", adr_content, re.DOTALL)
    if desc_match:
        module_data["description"] = desc_match.group(1).strip().split('\n')[0].strip()

    # Extract Criticality (specific to some ADRs or general search)
    crit_match = re.search(r"criticality:\s*[\"']?([A-Z_]+)[\"']?", adr_content, re.IGNORECASE)
    if crit_match:
        module_data["criticality"] = crit_match.group(1).upper()
    else: # Search general document for "criticality" keyword
        general_crit_match = re.search(r"(?:criticality|criticalidad):\s*(CRITICAL|HIGH|MEDIUM|LOW)", adr_content, re.IGNORECASE)
        if general_crit_match:
            module_data["criticality"] = general_crit_match.group(2).upper()
        # For MES it's medium, for others it could be derived from other parts.
        if module_name == "mes": module_data["criticality"] = "MEDIUM"
        if module_name in ["iam", "party", "product", "pricing", "sales"]: module_data["criticality"] = "CRITICAL"


    # Extract implementation_phase (from ADR-007 or within ADRs)
    phase_match = re.search(r"implementation_phase:\s*[\"']?(.+?)[\"']?", adr_content, re.IGNORECASE)
    if phase_match:
        module_data["implementation_phase"] = phase_match.group(1).strip()
    else: # Try to derive from ADR-007
        if module_name == "iam": module_data["implementation_phase"] = "0 (Before MVP)"
        elif module_name in ["party", "product", "pricing"]: module_data["implementation_phase"] = "1"
        elif module_name == "sales": module_data["implementation_phase"] = "2"
        elif module_name == "mes": module_data["implementation_phase"] = "3"

    # Extract Key Entities
    entities_match = re.search(r"(?:Key Entities|Entidades de Dominio Clave|Entidades Principales):\s*\n((?:- ['`].+?['`]\s*\n)+)", adr_content)
    if entities_match:
        for entity in entities_match.group(1).strip().split('\n'):
            module_data["key_entities"].append(entity.replace("- `", "").replace("`", "").strip())
    
    # Extract dependencies_on and depended_on_by (often in ADRs or diagrams)
    deps_on_match = re.search(r"dependencies_on:\s*\n((?:- ['`]?(.+?)['`]?\s*\n)+)", adr_content)
    if deps_on_match:
        for dep in deps_on_match.group(1).strip().split('\n'):
            module_data["dependencies_on"].append(dep.replace("- ", "").strip())

    depended_on_by_match = re.search(r"depended_on_by:\s*\n((?:- ['`]?(.+?)['`]?\s*\n)+)", adr_content)
    if depended_on_by_match:
        for dep in depended_on_by_match.group(1).strip().split('\n'):
            module_data["depended_on_by"].append(dep.replace("- ", "").strip())

    # Extract Use Cases (high-level from ADRs or linked module spec)
    use_cases_match = re.search(r"(?:use_cases|Casos de Uso):\s*\n((?:- ['`]?CU-\S+?['`]?\s*\n)+)", adr_content)
    if use_cases_match:
        for uc in use_cases_match.group(1).strip().split('\n'):
            module_data["use_cases"].append(uc.replace("- `", "").replace("`", "").strip())
    else:
        # General match for any list under "Casos de Uso" section
        uc_section_match = re.search(r"## \d+\.?\s*Casos de Uso(?:\s*\(Capa de Aplicación\))?\s*\n\n((?:- .+?\s*\n)+)", adr_content)
        if uc_section_match:
            for uc in uc_section_match.group(1).strip().split('\n'):
                if uc.startswith("- "):
                    module_data["use_cases"].append(uc.replace("- ", "").strip())


    # Extract API Endpoints (high-level from ADRs or linked API contracts)
    endpoints_match = re.search(r"(?:api_endpoints|Endpoints):\s*\n((?:- ['`]?/\S+?['`]?\s*\n)+)", adr_content)
    if endpoints_match:
        for ep in endpoints_match.group(1).strip().split('\n'):
            module_data["api_endpoints"].append(ep.replace("- `", "").replace("`", "").strip())
    else:
        ep_section_match = re.search(r"## \d+\.?\s*Contratos de API(?:\s*-\s*Módulo [A-Z]+)?\s*\n((?:- \*\*Endpoint:\*\* `\S+`\s*\n)+)", adr_content)
        if ep_section_match:
            for ep in ep_section_match.group(1).strip().split('\n'):
                ep_val_match = re.search(r"\*\*Endpoint:\*\* `(\S+)`", ep)
                if ep_val_match:
                    module_data["api_endpoints"].append(ep_val_match.group(1).strip())

    # Special handling for MES dependencies based on ADR-018 and ADR-019
    if module_name == "mes":
        module_data["dependencies_on"].extend(["Sales (production order source from ERP Core)", "Product (what to produce from ERP Core)", "Party (Client references from ERP Core)", "IAM (Inspector/User references from ERP Core)"])
        module_data["depended_on_by"].extend(["Inventory (stock completion in ERP Core)", "Reports (production metrics in ERP Core)"])
        # Remove duplicates
        module_data["dependencies_on"] = list(dict.fromkeys(module_data["dependencies_on"]))
        module_data["depended_on_by"] = list(dict.fromkeys(module_data["depended_on_by"]))

    return module_data

def update_bounded_contexts_yaml(adr_dir, target_yaml_path):
    module_adrs = {
        "iam": "ADR-014-iam-module-architecture.md",
        "party": "ADR-012-arquitectura-modulo-party.md",
        "product": "ADR-015-product-module-architecture.md",
        "pricing": "ADR-016-pricing-module-architecture.md",
        "sales": "ADR-017-sales-module-architecture.md",
        "mes": "ADR-018-mes-module-architecture.md"
    }

    core_contexts = {}
    for module_name, adr_filename in module_adrs.items():
        adr_path = os.path.join(adr_dir, adr_filename)
        if os.path.exists(adr_path):
            with open(adr_path, 'r', encoding='utf-8') as f:
                adr_content = f.read()
            core_contexts[module_name] = extract_from_adr(adr_content, module_name)
        else:
            print(f"Warning: ADR file not found for {module_name}: {adr_path}")

    # Load existing bounded-contexts.yaml template
    if os.path.exists(target_yaml_path):
        with open(target_yaml_path, 'r', encoding='utf-8') as f:
            data = yaml.safe_load(f)
    else:
        data = {}

    # Update the core_contexts section
    data["core_contexts"] = core_contexts

    # Update metadata
    if "metadata" not in data:
        data["metadata"] = {}
    data["metadata"]["last_updated"] = "2026-02-14" # Today's date (formatted according to the user's locale).

    # Save the updated YAML
    with open(target_yaml_path, 'w', encoding='utf-8') as f:
        yaml.safe_dump(data, f, sort_keys=False, default_flow_style=False, indent=2)

    print(f"Successfully updated {target_yaml_path}")


def update_architecture_yaml(adr_dir, architecture_vision_path, target_yaml_path):
    architecture_data = {
        "clean_architecture_layers": {
            "title": "🏗️ Clean Architecture Layers",
            "domain_layer": {},
            "application_layer": {},
            "infrastructure_layer": {},
            "interfaces_layer": {}
        },
        "ddd_strategy": {},
        "repository_pattern": {},
        "value_objects": {},
        "testing_strategy": {},
        "error_handling": {}
    }

    # Extract from ADR-002-clean-architecture-ddd-adoption.md
    adr002_path = os.path.join(adr_dir, "ADR-002-clean-architecture-ddd-adoption.md")
    if os.path.exists(adr002_path):
        with open(adr002_path, 'r', encoding='utf-8') as f:
            adr002_content = f.read()

        # Clean Architecture Layers (using existing data from template and enriching)
        architecture_data["clean_architecture_layers"]["domain_layer"] = {
            "name": "Domain (Core Business Logic)",
            "location": "apps/tramatex-api/internal/domain", # Use concrete backend app name
            "purpose": "Pure business logic, entities, value objects, domain services",
            "rules": [
                "NO external dependencies (no frameworks, no GORM, no HTTP)",
                "Pure Go: only standard library concepts",
                "Interfaces for external concerns (repositories, services)",
                "All business rules encapsulated here",
                "Completely testeable in isolation",
                "Typed errors (not strings)"
            ],
            "contains": [
                "Entities: Domain objects with identity",
                "Value Objects: Immutable domain concepts",
                "Domain Services: Cross-entity business logic",
                "Repositories: Interfaces (not implementations)",
                "Domain Events: (future enhancement)"
            ]
        }
        architecture_data["clean_architecture_layers"]["application_layer"] = {
            "name": "Application (Orchestration)",
            "location": "apps/tramatex-api/internal/application", # Use concrete backend app name
            "purpose": "Use cases, orchestration of domain and infrastructure",
            "rules": [
                "NO business logic (delegated to domain)",
                "Orchestrates domain and infrastructure",
                "Creates DTOs for layer boundaries",
                "Handles transactions and workflows",
                "Can inject dependencies (repositories, services)"
            ],
            "contains": [
                "Use Cases: Business workflows",
                "Application Services: Orchestrators",
                "DTOs: Layer boundary models",
                "Mappers: Domain ↔ Application ↔ Infrastructure"
            ]
        }
        architecture_data["clean_architecture_layers"]["infrastructure_layer"] = {
            "name": "Infrastructure (Technical Concerns)",
            "location": "apps/tramatex-api/internal/infrastructure", # Use concrete backend app name
            "purpose": "External concerns: database, APIs, libraries, technical plumbing",
            "rules": [
                "Implements domain interfaces (repositories)",
                "Contains framework-specific code (GORM, HTTP clients)",
                "NO business logic",
                "Substitutable: can swap implementations",
                "Technical, not strategic"
            ],
            "contains": [
                "Repository Implementations: GORM models + queries",
                "External Service Clients: 3rd party API wrappers",
                "Adapters: Technical translations",
                "Configuration: Database, logging, etc."
            ]
        }
        architecture_data["clean_architecture_layers"]["interfaces_layer"] = {
            "name": "Interfaces (Entry Points)",
            "location": "apps/tramatex-api/internal/interfaces", # Use concrete backend app name
            "purpose": "HTTP handlers, DTOs, request/response mapping",
            "rules": [
                "Controllers/Handlers: HTTP entry points",
                "Request/Response DTOs: HTTP contracts",
                "Input validation (basic, detailed in domain)",
                "Error translation to HTTP status codes",
                "No business logic"
            ],
            "contains": [
                "HTTP Handlers: Route handlers",
                "DTOs: HTTP request/response models",
                "Middleware: Logging, auth, CORS",
                "Error Handlers: HTTP error translation"
            ]
        }
        
        # DDD Strategy
        architecture_data["ddd_strategy"] = {
            "title": "📚 Domain-Driven Design Strategy",
            "asymmetric_rigor": {
                "concept": "Apply strict architecture where it adds value, relax where it's mechanical",
                "strict_rigor_applied_to": [
                    "Pricing (economically critical core engine)",
                    "Party (foundation for all operations)",
                    "Product (base for pricing calculations)",
                    "Sales (main business flow)",
                    "Authentication (security-critical)"
                ],
                "relaxed_rigor_applied_to": [
                    "Simple CRUDs of lookup data",
                    "Data migrations and utilities",
                    "Admin panels (post-MVP)",
                    "Simple reporting (no complex logic)"
                ],
                "principle": "Rule: If module has business logic → STRICT\nIf module is technical plumbing → PRAGMATIC"
            }
        }
    
    # Extract from ADR-011-testing-coverage-strategy.md
    adr011_path = os.path.join(adr_dir, "ADR-011-testing-coverage-strategy.md")
    if os.path.exists(adr011_path):
        with open(adr011_path, 'r', encoding='utf-8') as f:
            adr011_content = f.read()
        
        # Testing Strategy
        unit_tests_match = re.search(r"## \d+\.?\s*Decisión Adoptada.*?Pirámide de Testing: Se favorece una base sólida de tests unitarios rápidos, complementados por menos tests de integración y un número muy selectivo de tests E2E\. La proporción recomendada como guía es \*\*(.+?)\*\*\.", adr011_content, re.DOTALL)
        if unit_tests_match:
            architecture_data["testing_strategy"]["title"] = "🧪 Testing Strategy in TramaTex"
            architecture_data["testing_strategy"]["pyramid_principle"] = f"Pyramid of Testing: {unit_tests_match.group(1).strip()} Unitarios, Integración, E2E"
            
        coverage_table_match = re.search(r"#### MVP \(mínimos obligatorios\).*?\| Módulo \| Cobertura Mínima \| Criticidad \| Justificación \|\n\|---+\|---+\|---+\|---+\|\n((?:\|.+?\|\n)+)", adr011_content, re.DOTALL)
        if coverage_table_match:
            mvp_coverage_str = coverage_table_match.group(1).strip().split('\n')
            mvp_coverage = {}
            for line in mvp_coverage_str:
                parts = [p.strip() for p in line.split('|') if p.strip()]
                if len(parts) == 4:
                    module, cov, crit, just = parts
                    mvp_coverage[module.lower().replace(" ", "_")] = {"coverage": cov, "criticality": crit, "justification": just}
            architecture_data["testing_strategy"]["mvp_coverage_targets"] = mvp_coverage
        
        # Add testing philosophy, tools, etc. (dummy for now)
        architecture_data["testing_strategy"]["unit_tests"] = {
            "location": "In domain, application layers",
            "scope": "Individual entities, value objects, services",
            "tools": "Go testing + testify/assert"
        }
        architecture_data["testing_strategy"]["integration_tests"] = {
            "location": "In tests/integration/",
            "scope": "End-to-end use case validation with real database",
            "focus": [
                "Complete workflows (login → create order → pricing)",
                "Repository interactions",
                "Transaction boundaries"
            ],
            "tools": "Go testing + database fixtures"
        }
        architecture_data["testing_strategy"]["no_ui_tests_in_<backend_app_name>"] = {
            "rule": "tramatex-api tests domain and use cases, NOT HTTP handlers", # Concrete name
            "reason": "HTTP concerns tested separately or in frontend"
        }
    
    # Repository pattern (dummy for now)
    architecture_data["repository_pattern"] = {
        "title": "Repository Pattern (Interfaces in Domain)",
        "design": [
            "Repository interfaces defined in domain/",
            "Implementations in infrastructure/",
            "Domain depends on abstraction, not concretion",
            "Easy to test domain with mock repositories",
            "Easy to swap database implementations"
        ]
    }

    # Value Objects (dummy for now)
    architecture_data["value_objects"] = {
        "title": "Value Objects (Domain Primitives)",
        "concept": "Encapsulate domain concepts as immutable objects with validation",
        "examples": {
            "email": {"validation": "RFC 5322 format", "immutable": True, "example": "type Email struct { value string }"},
            "money": {"validation": "Positive value, decimal precision", "immutable": True, "example": "type Money struct { amount decimal.Decimal, currency string }"},
            "product_code": {"validation": "Format and uniqueness constraints", "immutable": True, "example": "type ProductCode struct { code string }"}
        }
    }

    # Error Handling (dummy for now)
    architecture_data["error_handling"] = {
        "title": "Error Handling in TramaTex",
        "principle": "Typed errors, not strings",
        "domain_errors": {
            "location": "apps/tramatex-api/pkg/errors/", # Concrete name
            "contains": [
                "Custom error types (InvalidEmail, ProductNotFound, etc.)",
                "Error codes (numeric or string constants)",
                "Error context (what went wrong, why, where)"
            ]
        },
        "application_errors": {
            "handling": "Catch domain errors, wrap for use case context"
        },
        "http_errors": {
            "translation": "Infrastructure layer translates to HTTP status codes"
        }
    }


    # Load existing architecture.yaml template
    if os.path.exists(target_yaml_path):
        with open(target_yaml_path, 'r', encoding='utf-8') as f:
            data = yaml.safe_load(f)
    else:
        data = {}

    # Update sections
    data.update(architecture_data)

    # Update metadata
    if "metadata" not in data:
        data["metadata"] = {}
    data["metadata"]["last_updated"] = "2026-02-14"

    # Save the updated YAML
    with open(target_yaml_path, 'w', encoding='utf-8') as f:
        yaml.safe_dump(data, f, sort_keys=False, default_flow_style=False, indent=2)
    
    print(f"Successfully updated {target_yaml_path}")

# Main execution
if __name__ == "__main__":
    adr_base_path = "docs/architecture/adrs"
    bounded_contexts_template_path = "agents/project/context/bounded-contexts.yaml"
    architecture_template_path = "agents/project/context/architecture.yaml"

    # Update bounded-contexts.yaml
    update_bounded_contexts_yaml(adr_base_path, bounded_contexts_template_path)

    # Update architecture.yaml
    update_architecture_yaml(adr_base_path, "docs/architecture/architecture-vision.md", architecture_template_path)