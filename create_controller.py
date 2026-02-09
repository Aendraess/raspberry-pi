import os
import json
import re
import tkinter as tk
from tkinter import ttk, messagebox
from typing import List, Optional
TESTING = True
TEMPLATE_FOLDER = os.path.join(os.path.dirname(__file__),'template')
CONTROLLER_FOLDER = os.path.join(os.path.dirname(__file__),'controllers')
DTO_FOLDER = os.path.join(os.path.dirname(__file__),'dtos')
MODELS_FOLDER = os.path.join(os.path.dirname(__file__),'models')
SERVICES_FOLDER = os.path.join(os.path.dirname(__file__),'services')
DATABASE_MIGRATIONS_FILE = os.path.join(os.path.dirname(__file__),'database','dbMigrations.go')
SERVER_FILE = os.path.join(os.path.dirname(__file__),'server','server.go')
CONTROLLER_TEMPLATE = os.path.join(TEMPLATE_FOLDER, 'controller.txt')
DTO_TEMPLATE = os.path.join(TEMPLATE_FOLDER, 'createDto.txt')
UPDATE_DTO_TEMPLATE = os.path.join(TEMPLATE_FOLDER, 'updateDto.txt')
MODEL_TEMPLATE = os.path.join(TEMPLATE_FOLDER, 'model.txt')
SERVICE_TEMPLATE = os.path.join(TEMPLATE_FOLDER, 'service.txt')
OUT_DIR = os.path.join(os.path.dirname(__file__),'out')

def decapitalize_first_letter(name: str) -> str:
        return name[0].lower() + name[1:]
    
class InputField:
    def __init__(self, name: str, type: str, existing_models: List[str]):
        self.name = name
        self.type = type
        self.existing_models = existing_models
        gorm_constraint = ",  gorm:\"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;\"" if self.type in self.existing_models else ""
        self.json = f"""`json:"{decapitalize_first_letter(self.name)}{gorm_constraint}"`"""
        print(*self.existing_models)
        print(["int", "float32", "string", "uint", "bool", *[model for model in self.existing_models]], self.name, self.type)
        assert self.type in ["int", "float32", "string", "uint", "bool", *[model for model in self.existing_models]]
        
        assert re.match(r'^[A-Za-z][A-Za-z0-9_]*$', self.name)
        assert self.name.lower() not in ["createdat", "updatedat", "deletedat", "basemodel"]
        assert self.name.lower() not in ["int", "float32", "string", "uint", "bool"]
        
        
    def __str__(self):
        return f"{self.name}: {self.type}"
    
    def __dict__(self):
        return {
            "name": self.name,
            "type": self.type,
            "json": self.json
        }
 
class ControllerCreator:    
    def __init__(self, model_name: str):
        self.model_name = model_name
        self.fields: List[InputField] = []
        self.existing_models = [model for model in self.collect_existing_models()]
        self.assert_files_exist()
    
    
    def collect_existing_models(self) -> List[str]:
        for file in os.listdir(MODELS_FOLDER):
            if file.endswith(".go"):
                with open(os.path.join(MODELS_FOLDER, file), 'r') as f:
                    content = f.read()
                    match = re.search(r'type\s+(\w+)\s+struct\s+{\n', content)
                    assert match, f"Model {file} does not have a type definition"
                    yield match.group(1)
        return []
    
    def assert_files_exist(self):
        os.makedirs(OUT_DIR, exist_ok=True)
        assert os.path.exists(TEMPLATE_FOLDER), f"Template folder not found at {TEMPLATE_FOLDER}"
        assert os.path.exists(CONTROLLER_FOLDER), f"Controller folder not found at {CONTROLLER_FOLDER}"
        assert os.path.exists(MODELS_FOLDER), f"Models folder not found at {MODELS_FOLDER}"
        assert os.path.exists(SERVICES_FOLDER), f"Services folder not found at {SERVICES_FOLDER}"
        assert os.path.exists(DATABASE_MIGRATIONS_FILE), f"Database migrations file not found at {DATABASE_MIGRATIONS_FILE}"
        assert os.path.exists(SERVER_FILE), f"Server file not found at {SERVER_FILE}"
        assert os.path.exists(CONTROLLER_TEMPLATE), f"Controller template file not found at {CONTROLLER_TEMPLATE}"
        assert os.path.exists(DTO_TEMPLATE), f"Create DTO template file not found at {DTO_TEMPLATE}"
        assert os.path.exists(UPDATE_DTO_TEMPLATE), f"Update DTO template file not found at {UPDATE_DTO_TEMPLATE}"
        assert os.path.exists(MODEL_TEMPLATE), f"Model template file not found at {MODEL_TEMPLATE}"
        assert os.path.exists(SERVICE_TEMPLATE), f"Service template file not found at {SERVICE_TEMPLATE}"
        assert os.path.exists(OUT_DIR), f"Out directory not found at {OUT_DIR}"
    
    def add_field(self, name: str, type: str) -> None:
        field = InputField(name, type, self.existing_models)
        assert not any(f.name == name for f in self.fields), f"Field name {name} already exists"
        assert field.name.lower() not in ["createdat", "updatedat", "deletedat", "basemodel"], f"Field name {name} is reserved"
        self.fields.append(field)
        if field.type in self.existing_models:
            assert f"{decapitalize_first_letter(field.type)}Id" not in [field.name for field in self.fields], f"Field name {name} is reserved"
            assert field.type in self.existing_models, f"Model {field.type} does not exist"
            self.fields.append(InputField(f"{decapitalize_first_letter(field.type)}Id", "uint", self.existing_models))
    
    def clear_fields(self) -> None:
        self.fields.clear()
    
    def generate(self) -> None:
        assert self.fields, "Please add at least one field"
        
        # Read templates
        with open(CONTROLLER_TEMPLATE, 'r') as f:
            controller_template = f.read()
        with open(DTO_TEMPLATE, 'r') as f:
            dto_template = f.read()
        with open(UPDATE_DTO_TEMPLATE, 'r') as f:
            update_dto_template = f.read()
        with open(MODEL_TEMPLATE, 'r') as f:
            model_template = f.read()
        with open(SERVICE_TEMPLATE, 'r') as f:
            service_template = f.read()
        
        self._generate_controller(controller_template)
        self._generate_dtos(dto_template, update_dto_template)
        self._generate_model(model_template)
        self._generate_service(service_template)
        self._update_database_migrations()
        self._update_server()
    
    def _generate_controller(self, template: str) -> None:
        
        controller_name = f"{self.model_name}Controller"
        controller_file = os.path.join(CONTROLLER_FOLDER, f"{controller_name}.go")
        assert not os.path.exists(controller_file), f"Controller file already exists at {controller_file}"
        
        pass
    
    def _generate_dtos(self, create_template: str, update_template: str) -> None:
        CreateDtoName = f"Create{self.model_name}Request"
        UpdateDtoName = f"Update{self.model_name}Request"
        assert create_template != None and create_template != ""
        assert update_template != None and update_template != ""
        longest_field_name = len(max(self.fields, key=lambda x: len(x.name)).name)
        longest_field_type = len(max(self.fields, key=lambda x: len(self.convert_type(x.type))).type)    
        create_template = create_template.replace("[UpdateDto]",UpdateDtoName)
        create_template = create_template.replace("    [Fields]",'\n'.join(["    "+str(field.name).ljust(longest_field_name) + f"   {self.convert_type(field.type).ljust(longest_field_type)}   {field.json}" for field in self.fields]))
        
        
        update_template = update_template.replace("[UpdateDto]",UpdateDtoName)
        update_template = update_template.replace("    [Fields]",'\n'.join(["    "+str(field.name).ljust(longest_field_name) + f"   {self.convert_type(field.type).ljust(longest_field_type)}   {field.json}" for field in self.fields]))
        
        updatepath = os.path.join(DTO_FOLDER if TESTING else TEMPLATE_FOLDER,f"{UpdateDtoName}.go")
        createpath = os.path.join(DTO_FOLDER  if TESTING else TEMPLATE_FOLDER,f"{CreateDtoName}.go")
        with open(updatepath,'w', encoding='utf-8') as f:
            f.write(update_template)
        pass
    
    def convert_type(self, type: str) -> str:
        if type == "int":
            return "int"
        elif type == "float32":
            return "float32"
        elif type == "string":
            return "string"
        elif type == "bool":
            return "bool"
        elif type == "uint":
            return "uint"
        else:
            return type
    
    def _generate_model(self, template: str) -> None:
        # Implementation for model generation
        filepath = os.path.join(OUT_DIR, decapitalize_first_letter(self.model_name) + ".go") if TESTING else os.path.join(MODELS_FOLDER, self.model_name + ".go")
        if not TESTING:
            assert not os.path.exists(filepath), f"Model file already exists at {filepath}"
        with open(filepath, 'w') as f:
            longest_field_name = len(max(self.fields, key=lambda x: len(x.name)).name)
            longest_field_type = len(max(self.fields, key=lambda x: len(self.convert_type(x.type))).type)
            new_template = template.replace("[EntityName]",self.model_name).replace("    [Fields]",'\n'.join(["    "+str(field.name).ljust(longest_field_name) + f"   {self.convert_type(field.type).ljust(longest_field_type)}   {field.json}" for field in self.fields]))
            
            f.write(new_template)
    
    def _generate_service(self, template: str) -> None:
        # Implementation for service generation
        pass
    
    def _update_database_migrations(self) -> None:
        # Implementation for database migrations update
        pass
    
    def _update_server(self) -> None:
        # Implementation for server update
        pass

class ControllerCreatorGUI:
    def __init__(self, root):
        self.root = root
        self.root.title("Controller Creator")
        
        self.controller_creator = ControllerCreator("MarketItem")
        self.setup_ui()
        
    def setup_ui(self):
        # Main frame
        main_frame = ttk.Frame(self.root, padding="10")
        main_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Model name entry
        ttk.Label(main_frame, text="Model Name:").grid(row=0, column=0, sticky=tk.W, pady=5)
        self.model_name_var = tk.StringVar(value="MarketItem")
        self.model_name_entry = ttk.Entry(main_frame, textvariable=self.model_name_var)
        self.model_name_entry.grid(row=0, column=1, sticky=(tk.W, tk.E), pady=5)
        
        # Fields frame
        fields_frame = ttk.LabelFrame(main_frame, text="Fields", padding="5")
        fields_frame.grid(row=1, column=0, columnspan=2, sticky=(tk.W, tk.E, tk.N, tk.S), pady=5)
        
        # Fields list
        self.fields_tree = ttk.Treeview(fields_frame, columns=("name", "type"), show="headings")
        self.fields_tree.heading("name", text="Name")
        self.fields_tree.heading("type", text="Type")
        self.fields_tree.grid(row=0, column=0, columnspan=2, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Scrollbar for fields list
        scrollbar = ttk.Scrollbar(fields_frame, orient=tk.VERTICAL, command=self.fields_tree.yview)
        scrollbar.grid(row=0, column=2, sticky=(tk.N, tk.S))
        self.fields_tree.configure(yscrollcommand=scrollbar.set)
        
        # Add field frame
        add_field_frame = ttk.Frame(fields_frame)
        add_field_frame.grid(row=1, column=0, columnspan=2, sticky=(tk.W, tk.E), pady=5)
        
        ttk.Label(add_field_frame, text="Name:").grid(row=0, column=0, padx=5)
        self.field_name_var = tk.StringVar()
        self.field_name_entry = ttk.Entry(add_field_frame, textvariable=self.field_name_var)
        self.field_name_entry.grid(row=0, column=1, padx=5)
        
        ttk.Label(add_field_frame, text="Type:").grid(row=0, column=2, padx=5)
        self.field_type_var = tk.StringVar()
        self.field_type_combo = ttk.Combobox(add_field_frame, textvariable=self.field_type_var, 
                                           values=["int", "float", "str", "bool"], state="readonly")
        self.field_type_combo.grid(row=0, column=3, padx=5)
        
        ttk.Button(add_field_frame, text="Add Field", command=self.add_field).grid(row=0, column=4, padx=5)
        
        # Buttons
        button_frame = ttk.Frame(main_frame)
        button_frame.grid(row=2, column=0, columnspan=2, pady=10)
        
        ttk.Button(button_frame, text="Generate", command=self.generate).grid(row=0, column=0, padx=5)
        ttk.Button(button_frame, text="Clear", command=self.clear_fields).grid(row=0, column=1, padx=5)
        
        # Configure grid weights
        main_frame.columnconfigure(1, weight=1)
        main_frame.rowconfigure(1, weight=1)
        fields_frame.columnconfigure(0, weight=1)
        fields_frame.rowconfigure(0, weight=1)
        
    def add_field(self):
        name = self.field_name_var.get().strip()
        type = self.field_type_var.get()
        
        if not name or not type:
            messagebox.showerror("Error", "Please enter both name and type")
            return
            
        try:
            self.controller_creator.add_field(name, type)
            self.fields_tree.insert("", tk.END, values=(name, type))
            
            # Clear inputs
            self.field_name_var.set("")
            self.field_type_var.set("")
            
        except (AssertionError, ValueError) as e:
            messagebox.showerror("Error", str(e))
    
    def clear_fields(self):
        self.controller_creator.clear_fields()
        for item in self.fields_tree.get_children():
            self.fields_tree.delete(item)
    
    def generate(self):
        try:
            self.controller_creator.generate()
            messagebox.showinfo("Success", f"Successfully generated controller for {self.model_name_var.get()}")
        except Exception as e:
            messagebox.showerror("Error", str(e))

def main():
    root = tk.Tk()
    app = ControllerCreatorGUI(root)
    root.mainloop()

if __name__ == "__main__":
    # main()
    c = ControllerCreator("LogbookEntry")
    c.add_field("message", "string")
    c.add_field("level", "string")
    c.add_field("category", "string")
    c.add_field("timestamp", "time.Time")
    c.generate()
