import { Component } from '@angular/core';
import { InventoryService } from '../services/inventory.services';

@Component({
  selector: 'app-product-list',
  templateUrl: './product-list.component.html'
})
export class ProductListComponent {
  newProduct = { name: '', price: 0, stock: 0 };

  constructor(private inventoryService: InventoryService) { }

  createProduct() {
    this.inventoryService.createProduct(this.newProduct).subscribe({
      next: () => alert('Produto criado!'),
      error: (err) => alert('Erro: ' + err.error)
    });
  }
}
