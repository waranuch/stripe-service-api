#!/usr/bin/env python3
"""
OpenAPI Specification Validator

This script validates the OpenAPI specification file to ensure it's properly formatted
and contains all required fields.
"""

import yaml
import json
import sys
import os
from typing import Dict, Any, List

def load_openapi_spec(file_path: str) -> Dict[str, Any]:
    """Load and parse the OpenAPI specification file."""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return yaml.safe_load(f)
    except FileNotFoundError:
        print(f"❌ Error: OpenAPI file not found: {file_path}")
        sys.exit(1)
    except yaml.YAMLError as e:
        print(f"❌ Error: Invalid YAML syntax in {file_path}: {e}")
        sys.exit(1)

def validate_basic_structure(spec: Dict[str, Any]) -> List[str]:
    """Validate basic OpenAPI structure."""
    errors = []
    
    # Check required root fields
    required_fields = ['openapi', 'info', 'paths']
    for field in required_fields:
        if field not in spec:
            errors.append(f"Missing required field: {field}")
    
    # Check OpenAPI version
    if 'openapi' in spec:
        version = spec['openapi']
        if not version.startswith('3.'):
            errors.append(f"Unsupported OpenAPI version: {version}")
    
    # Check info section
    if 'info' in spec:
        info = spec['info']
        required_info_fields = ['title', 'version']
        for field in required_info_fields:
            if field not in info:
                errors.append(f"Missing required info field: {field}")
    
    return errors

def validate_paths(spec: Dict[str, Any]) -> List[str]:
    """Validate API paths and operations."""
    errors = []
    
    if 'paths' not in spec:
        return errors
    
    paths = spec['paths']
    expected_paths = [
        '/health',
        '/customers',
        '/customers/{id}',
        '/payment-intents',
        '/payment-intents/{id}/confirm',
        '/products',
        '/prices',
        '/subscriptions',
        '/subscriptions/{id}'
    ]
    
    # Check if all expected paths exist
    for path in expected_paths:
        if path not in paths:
            errors.append(f"Missing API path: {path}")
    
    # Validate each path
    for path, path_item in paths.items():
        if not isinstance(path_item, dict):
            errors.append(f"Invalid path item for {path}")
            continue
        
        # Check operations
        for method, operation in path_item.items():
            if method.upper() not in ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']:
                continue
            
            if not isinstance(operation, dict):
                errors.append(f"Invalid operation for {method.upper()} {path}")
                continue
            
            # Check required operation fields
            required_op_fields = ['summary', 'operationId', 'tags', 'responses']
            for field in required_op_fields:
                if field not in operation:
                    errors.append(f"Missing {field} in {method.upper()} {path}")
    
    return errors

def validate_components(spec: Dict[str, Any]) -> List[str]:
    """Validate components section."""
    errors = []
    
    if 'components' not in spec:
        errors.append("Missing components section")
        return errors
    
    components = spec['components']
    
    # Check schemas
    if 'schemas' not in components:
        errors.append("Missing schemas in components")
        return errors
    
    schemas = components['schemas']
    expected_schemas = [
        'Customer',
        'CreateCustomerRequest',
        'ListCustomersResponse',
        'PaymentIntent',
        'CreatePaymentIntentRequest',
        'ConfirmPaymentIntentRequest',
        'Product',
        'CreateProductRequest',
        'Price',
        'CreatePriceRequest',
        'Subscription',
        'CreateSubscriptionRequest',
        'Error'
    ]
    
    for schema_name in expected_schemas:
        if schema_name not in schemas:
            errors.append(f"Missing schema: {schema_name}")
    
    return errors

def validate_responses(spec: Dict[str, Any]) -> List[str]:
    """Validate response definitions."""
    errors = []
    
    if 'components' not in spec or 'responses' not in spec['components']:
        errors.append("Missing responses in components")
        return errors
    
    responses = spec['components']['responses']
    expected_responses = ['BadRequest', 'NotFound', 'InternalServerError']
    
    for response_name in expected_responses:
        if response_name not in responses:
            errors.append(f"Missing response: {response_name}")
    
    return errors

def count_endpoints(spec: Dict[str, Any]) -> int:
    """Count total number of endpoints."""
    count = 0
    if 'paths' in spec:
        for path, path_item in spec['paths'].items():
            for method in path_item.keys():
                if method.upper() in ['GET', 'POST', 'PUT', 'DELETE', 'PATCH']:
                    count += 1
    return count

def print_statistics(spec: Dict[str, Any]):
    """Print specification statistics."""
    print("\n📊 OpenAPI Specification Statistics:")
    print(f"   OpenAPI Version: {spec.get('openapi', 'Unknown')}")
    print(f"   API Title: {spec.get('info', {}).get('title', 'Unknown')}")
    print(f"   API Version: {spec.get('info', {}).get('version', 'Unknown')}")
    print(f"   Total Endpoints: {count_endpoints(spec)}")
    
    if 'paths' in spec:
        print(f"   Total Paths: {len(spec['paths'])}")
    
    if 'components' in spec and 'schemas' in spec['components']:
        print(f"   Total Schemas: {len(spec['components']['schemas'])}")
    
    if 'servers' in spec:
        print(f"   Servers: {len(spec['servers'])}")
        for i, server in enumerate(spec['servers']):
            print(f"     {i+1}. {server.get('url', 'Unknown')} - {server.get('description', 'No description')}")

def main():
    """Main validation function."""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    openapi_file = os.path.join(script_dir, '..', 'openapi.yaml')
    
    print("🔍 Validating OpenAPI Specification...")
    print(f"📄 File: {openapi_file}")
    
    # Load specification
    spec = load_openapi_spec(openapi_file)
    
    # Run validations
    all_errors = []
    all_errors.extend(validate_basic_structure(spec))
    all_errors.extend(validate_paths(spec))
    all_errors.extend(validate_components(spec))
    all_errors.extend(validate_responses(spec))
    
    # Print results
    if all_errors:
        print("\n❌ Validation Errors:")
        for error in all_errors:
            print(f"   • {error}")
        print(f"\n📊 Total Errors: {len(all_errors)}")
        sys.exit(1)
    else:
        print("\n✅ OpenAPI Specification is valid!")
        print_statistics(spec)
        print("\n🎉 All validations passed!")

if __name__ == "__main__":
    main() 