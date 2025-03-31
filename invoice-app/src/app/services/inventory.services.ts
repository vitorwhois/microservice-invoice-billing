import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class InventoryService {
    private apiUrl = 'http://localhost:8080';

    constructor(private http: HttpClient) { }

    createProduct(product: any) {
        return this.http.post(`${this.apiUrl}/products`, product);
    }

    getProduct(id: number) {
        return this.http.get(`${this.apiUrl}/products/${id}`);
    }
}
