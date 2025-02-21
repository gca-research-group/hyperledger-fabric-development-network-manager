import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarComponent } from '@app/components/sidebar';
import { ToolbarComponent } from './components/toolbar';
import { Sidebar } from './models';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  imports: [RouterOutlet, SidebarComponent, ToolbarComponent],
})
export class AppComponent {
  menus: Sidebar[] = [
    { label: 'home', icon: 'home', url: '' },
    { label: 'configuration-files', icon: 'link', url: 'configuration-files' },
    { label: 'settings', icon: 'settings', url: '' },
  ];
}
