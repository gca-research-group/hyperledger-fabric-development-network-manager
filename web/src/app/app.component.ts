import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';

import { SidebarComponent } from '@app/components/sidebar';

import { ToolbarComponent } from './components/toolbar';
import { Sidebar } from './models';
import { CurrentUserService } from './services/current-user';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  imports: [RouterOutlet, SidebarComponent, ToolbarComponent],
})
export class AppComponent {
  private currentUserService = inject(CurrentUserService);
  currentUser$ = this.currentUserService.currentUser$;

  menus: Sidebar[] = [
    { label: 'home', icon: 'home', url: '' },
    { label: 'orderer', icon: 'swap_vert', url: 'orderer' },
    { label: 'peer', icon: 'network_node', url: 'peer' },
    { label: 'channel', icon: 'hub', url: 'channel' },
  ];
}
