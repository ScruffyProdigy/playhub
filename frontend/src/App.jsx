import CompleteSignInPage from './components/auth/CompleteSignInPage'
import ReturnPage from './components/auth/ReturnPage'
import AuthPanel from './components/auth/AuthPanel'
import IntentBanner from './components/games/IntentBanner'
import GameLobby from './components/games/GameLobby'
import CreateRoomPanel from './components/rooms/CreateRoomPanel'
import { ActiveRoomProvider, useActiveRoom } from './components/rooms/ActiveRoomProvider'
import AppDock from './components/rooms/AppDock'
import RoomPanel from './components/rooms/RoomPanel'
import RoomSheet from './components/rooms/RoomSheet'
import { AuthProvider, useAuth } from './components/auth/AuthProvider'
import { useActiveIntent } from './components/games/useActiveIntent'
import { APP_NAME, APP_TAGLINE } from './lib/brand'
import { parseRoomInviteCode } from './lib/rooms'
import { MOBILE_ROOM_QUERY, useMediaQuery } from './lib/useMediaQuery'
import { usePathname } from './lib/usePathname'
import { useEffect } from 'react'
import './App.css'

function CatalogPage() {
  const { user, loading: authLoading } = useAuth()
  const { activeIntent, activeTableSeat, busy, refresh, notifyQueueJoined, handleLeave } = useActiveIntent()

  return (
    <main className="app-shell app-shell--catalog">
      {!authLoading && user ? (
        <IntentBanner
          activeIntent={activeIntent}
          activeTableSeat={activeTableSeat}
          busy={busy}
          onLeave={handleLeave}
        />
      ) : null}

      <header className="app-header">
        <h1>{APP_NAME}</h1>
        <p className="tagline">{APP_TAGLINE}</p>
      </header>

      <AuthPanel />
      {!authLoading && user ? <CreateRoomPanel /> : null}
      <GameLobby
        activeIntent={activeIntent}
        activeTableSeat={activeTableSeat}
        onQueueChange={refresh}
        onQueueJoined={notifyQueueJoined}
        onTableChange={refresh}
      />
    </main>
  )
}

function MainLayout() {
  const pathname = usePathname()
  const inviteCode = parseRoomInviteCode(pathname)
  const { user } = useAuth()
  const { room, roomOpen, dismissRoom, openRoom, unreadCount, hasRoomMembership, markRead } = useActiveRoom()
  const isMobile = useMediaQuery(MOBILE_ROOM_QUERY)
  const inRoomContext = Boolean(room || inviteCode || hasRoomMembership)
  // Desktop keeps the room panel visible whenever the user belongs to a room; mobile toggles via dock/sheet.
  const showDesktopRoom = Boolean(room && user && !isMobile)
  const showMobileSheet = Boolean(user && isMobile && roomOpen)
  const showDock = Boolean(user && isMobile && inRoomContext)

  useEffect(() => {
    if (showDesktopRoom) {
      markRead()
    }
  }, [showDesktopRoom, markRead, room?.id])

  useEffect(() => {
    const root = document.getElementById('root')
    root?.classList.toggle('app-root--split', showDesktopRoom)
    root?.classList.toggle('app-root--dock', showDock)
    return () => {
      root?.classList.remove('app-root--split', 'app-root--dock')
    }
  }, [showDesktopRoom, showDock])

  return (
    <>
      <div className={`app-layout ${showDesktopRoom ? 'app-layout--split' : ''}`}>
        <div className="app-layout__catalog">
          <CatalogPage />
        </div>
        {showDesktopRoom ? (
          <aside className="app-layout__room panel-card" aria-label="Room chat">
            <RoomPanel compact />
          </aside>
        ) : null}
      </div>

      {showMobileSheet ? (
        <RoomSheet open={roomOpen} onDismiss={dismissRoom}>
          <RoomPanel compact />
        </RoomSheet>
      ) : null}

      {showDock ? (
        <AppDock
          unreadCount={unreadCount}
          roomOpen={roomOpen}
          onOpenCatalog={dismissRoom}
          onOpenRoom={openRoom}
        />
      ) : null}
    </>
  )
}

function MainShell() {
  const pathname = usePathname()
  const inviteCode = parseRoomInviteCode(pathname)

  return (
    <ActiveRoomProvider pendingInviteCode={inviteCode}>
      <MainLayout />
    </ActiveRoomProvider>
  )
}

function App() {
  const pathname = usePathname()
  const isMainRoute = pathname === '/' || parseRoomInviteCode(pathname)

  return (
    <AuthProvider>
      {pathname.startsWith('/auth/complete') ? (
        <CompleteSignInPage />
      ) : pathname.startsWith('/return') ? (
        <ReturnPage />
      ) : isMainRoute ? (
        <MainShell />
      ) : (
        <main className="app-shell auth-page">
          <h1>Page not found</h1>
          <a className="auth-link" href="/">
            Back to {APP_NAME}
          </a>
        </main>
      )}
    </AuthProvider>
  )
}

export default App
