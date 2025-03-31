import { ComponentFixture, TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { PeerSelectorComponent } from './peer-selector.component';

describe('PeerSelectorComponent', () => {
  let component: PeerSelectorComponent;
  let fixture: ComponentFixture<PeerSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PeerSelectorComponent],
      providers: [appConfig.providers],
    }).compileComponents();

    fixture = TestBed.createComponent(PeerSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
