import { Component } from '@angular/core';
import { InventoryService } from '../services/inventory.services';

@Component({
  selector: 'app-product-list',
  templateUrl: './product-list.component.html'
})
export class ProductListComponent {
  newProduct = { name: '', price: 0, stock: 0 };
  products: any[] = [];

  constructor(private inventoryService: InventoryService) { }

  ngOnInit(): void {
    this.getAllProducts();
  }


  createProduct() {
    if (!this.newProduct.name || this.newProduct.price <= 0) {
      alert('Preencha os dados do produto corretamente.');
      return;
    }
    this.inventoryService.createProduct(this.newProduct).subscribe({
      next: () => {
        alert('Produto criado!');
        this.newProduct = { name: '', price: 0, stock: 0 };
        this.getAllProducts();
      },
      error: (err) => alert('Erro: ' + err.error)
    });
  }
  getAllProducts() {
    this.inventoryService.getAllProducts().subscribe({
      next: (data) => this.products = data,
      error: (err) => alert('Erro ao carregar os produtos: ' + err.error)
    });
  }
}
