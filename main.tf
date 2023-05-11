provider "azurerm" {
  skip_provider_registration = "true"
  features {}
}

resource "azurerm_resource_group" "app-rg" {
  name     = "MyDemoAppRG"
  location = "eastus2"
}

resource "azurerm_service_plan" "myapp" {
  name                = "myapp-serviceplan"
  location            = azurerm_resource_group.app-rg.location
  resource_group_name = azurerm_resource_group.app-rg.name

  os_type  = "Linux"
  sku_name = "F1"
}

resource "azurerm_linux_web_app" "myapp" {
  name                = "myapp-demo-pov"
  resource_group_name = azurerm_resource_group.app-rg.name
  location            = azurerm_resource_group.app-rg.location
  service_plan_id     = azurerm_service_plan.myapp.id


  site_config {
    application_stack {
      go_version = "1.19"
    }
    always_on = false
  }

  identity {
    type = "SystemAssigned"
  }

}

# resource "azurerm_app_service_source_control" "myapp" {
#   app_id = azurerm_linux_web_app.myapp.id
#   repo_url = ""
#   branch = "master"
# }

output "creds" {
  value = azurerm_linux_web_app.myapp.identity
}